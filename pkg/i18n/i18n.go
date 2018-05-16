package i18n

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chai2010/gettext-go/gettext"
	"github.com/golang/glog"
)

var knownTranslations = map[string][]string{
	"cmdctl": {
		"default",
		"en_US",
		"fr_FR",
		"zh_CN",
		"ja_JP",
		"zh_TW",
		"it_IT",
		"de_DE",
	},
	// only used for unit tests.
	"test": {
		"default",
		"en_US",
	},
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	//"translations/test/en_US/LC_MESSAGES/grape.po": translationsTestEn_usLc_messagesK8sPo,
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

func loadSystemLanguage() string {
	langStr := os.Getenv("LANG")
	if langStr == "" {
		glog.V(3).Infof("Couldn't find the LANG environment variable, defaulting to en_US")
		return "default"
	}
	pieces := strings.Split(langStr, ".")
	if len(pieces) != 2 {
		glog.V(3).Infof("Unexpected system language (%s), defaulting to en_US", langStr)
		return "default"
	}
	return pieces[0]
}

func findLanguage(root string, getLanguageFn func() string) string {
	langStr := getLanguageFn()

	translations := knownTranslations[root]
	if translations != nil {
		for ix := range translations {
			if translations[ix] == langStr {
				return langStr
			}
		}
	}
	glog.V(3).Infof("Couldn't find translations for %s, using default", langStr)
	return "default"
}

// LoadTranslations loads translation files. getLanguageFn should return a language
// string (e.g. 'en-US'). If getLanguageFn is nil, then the loadSystemLanguage function
// is used, which uses the 'LANG' environment variable.
func LoadTranslations(root string, getLanguageFn func() string) error {
	if getLanguageFn == nil {
		getLanguageFn = loadSystemLanguage
	}

	langStr := findLanguage(root, getLanguageFn)
	translationFiles := []string{
		fmt.Sprintf("%s/%s/LC_MESSAGES/grape.po", root, langStr),
		fmt.Sprintf("%s/%s/LC_MESSAGES/grape.mo", root, langStr),
	}

	glog.V(3).Infof("Setting language to %s", langStr)
	// TODO: list the directory and load all files.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Make sure to check the error on Close.
	for _, file := range translationFiles {
		filename := "translations/" + file
		f, err := w.Create(file)
		if err != nil {
			return err
		}

		data, err := Asset(filename)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return nil
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	gettext.BindTextdomain("grape", root+".zip", buf.Bytes())
	gettext.Textdomain("grape")
	gettext.SetLocale(langStr)
	return nil
}

// T translates a string, possibly substituting arguments into it along
// the way. If len(args) is > 0, args1 is assumed to be the plural value
// and plural translation is used.
func T(defaultValue string, args ...int) string {
	if len(args) == 0 {
		return gettext.PGettext("", defaultValue)
	}
	return fmt.Sprintf(gettext.PNGettext("", defaultValue, defaultValue+".plural", args[0]),
		args[0])
}

// Errorf produces an error with a translated error string.
// Substitution is performed via the `T` function above, following
// the same rules.
func Errorf(defaultValue string, args ...int) error {
	return errors.New(T(defaultValue, args...))
}
