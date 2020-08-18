// crate by d1y<chenhonzhou@gmail.com>

package macapp

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// AppConfig app config
type AppConfig struct {
	AppPath string // the app path (z-index: 10)
	AppName string // the app name
	AppRoot bool   // if `true` find `/Application` else find `pwd/` (z-index: 20)
}

// AppRes ff
type AppRes struct {
	Conf AppConfig // config
}

// CreateRootAppPath create root app just like `/Application/yoxi.app`
func CreateRootAppPath(appname string) string {
	// osx app root path
	var root = `/Application`
	return filepath.Join(root, fmt.Sprintf(`%v.app`, appname))
}

// New new app
func New(conf AppConfig) AppRes {
	var n, r, p = conf.AppName, conf.AppRoot, conf.AppPath
	var f = fmt.Sprintf(`%v.app`, n)
	var P = path.Join(curr(), f)
	if len(p) >= 1 {
		var x = path.Join(p, f)
		if ensureDir(x) {
			P = x
		}
	}
	if r {
		P = CreateRootAppPath(n)
	}
	conf.AppPath = P
	return AppRes{
		Conf: conf,
	}
}

// CreateFolder create `name.app` folder
func (res AppRes) CreateFolder() bool {
	var Fpath = res.GetPath()
	return ensureDir(Fpath)
}

// CreateInitInfoPlist create default `info.plist` file data
// man link: https://www.jianshu.com/p/e8c686a00027
func CreateInitInfoPlist(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDisplayName</key>
  <string>%v</string>
	<key>CFBundleIconFile</key>
	<string>%v</string>
</dict>
</plist>`, name, name)
}

// SetBinFile set `bin` file, just like golang run `func main(){}`
// copy `input` bin file to `name.app/Contents/MacOS/name`
func (res AppRes) SetBinFile(input string) (int64, error) {
	var dist = res.GetBinPath()
	return copy(input, dist)
}

// SetIconByIcns set icon by `icns` format
func (res AppRes) SetIconByIcns(input string) (int64, error) {
	var r = res.GetIconPath()
	return copy(input, r)
}

// SetIcon set icon, support `png`/`jpg`/`gif`/`..` fotmat
// if icon format is `.icns`, please use [SetIconByIcns] method
// use imagemagick format to `.incs`, before install this
// website: https://imagemagick.org
// github: https://github.com/ImageMagick
func (res AppRes) SetIcon(input string) ([]byte, error) {
	var r = res.GetIconPath()
	cache := exec.Command(`magick`, input, r)
	return cache.CombinedOutput()
}

// CreateAppContentFolder create `name.app` content folder
func (res AppRes) CreateAppContentFolder() {
	var Fpath = res.GetPath()
	var Fname = res.GetName()

	// Contents folder
	var CtxPath = path.Join(Fpath, "Contents")

	// the folder is `cli` bin file
	// if create `yoxi.app/Contents/MacOS`
	// open `yoxi.app`, first run `yoxi.app/Contents/MacOS/yoxi`
	var RunCliPath = path.Join(CtxPath, "MacOS")
	ensureDir(RunCliPath)

	// this is `name.app` resoures
	// if set `name.app` icon, add `name.app/Contents/Resources/name.icns`
	var ResPath = path.Join(CtxPath, "Resources")
	ensureDir(ResPath)

	// create `info.plist`, add sample content
	var infoPlistPath = path.Join(CtxPath, "info.plist")
	var infoPlistData = CreateInitInfoPlist(Fname)
	ioutil.WriteFile(infoPlistPath, []byte(infoPlistData), 0755)

}

// GetPath get the path
// `New` is auto set root path
func (res AppRes) GetPath() string {
	return res.Conf.AppPath
}

// GetName get the name
func (res AppRes) GetName() string {
	return res.Conf.AppName
}

// GetIconPath get icon path
func (res AppRes) GetIconPath() string {
	var Fpath = res.GetPath()
	var Fname = res.GetName()
	var r = path.Join(Fpath, fmt.Sprintf("./Contents/Resources/%v.icns", Fname))
	return r
}

// GetBinPath get bin file path
func (res AppRes) GetBinPath() string {
	var Fpath = res.GetPath()
	var Fname = res.GetName()
	var dist = path.Join(Fpath, fmt.Sprintf("./Contents/MacOS/%v", Fname))
	return dist
}

// Create crate app
func Create(conf AppConfig) AppRes {
	res := New(conf)
	res.CreateFolder()
	res.CreateAppContentFolder()
	return res
}

func curr() string {
	d, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	return d
}

// ensureDir  auto create files and folders
func ensureDir(fileName string) bool {
	if _, serr := os.Stat(fileName); serr != nil {
		merr := os.MkdirAll(fileName, 0755)
		if merr != nil {
			return false
		}
	}
	return true
}

// copy file
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
