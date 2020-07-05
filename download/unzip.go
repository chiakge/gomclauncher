package download

import (
	"archive/zip"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/xmdhs/gomclauncher/launcher"
)

func (l Libraries) Unzip(typee string, i int) error {
	e, done, ch := creatch(len(l.librarie.Libraries), i)
	natives := make([]string, 0)
	m := sync.Mutex{}
	go func() {
		for _, v := range l.librarie.Libraries {
			v := v
			path, sha1, url := swichnatives(v)
			path = launcher.Minecraft + `/libraries/` + path
			if url == "" {
				done <- true
				continue
			}
			if ifallow(v) {
				m.Lock()
				natives = append(natives, path)
				m.Unlock()
			}
			if ifallow(v) && !ver(path, sha1) {
				if path != "" {
					d := downinfo{
						typee: typee,
						url:   url,
						path:  path,
						e:     e,
						Sha1:  sha1,
						done:  done,
						ch:    ch,
					}
					ch <- true
					go d.down()
				}
			} else {
				done <- true
			}

		}
	}()
	n := 0
	for {
		select {
		case <-done:
			n++
			if n == len(l.librarie.Libraries) {
				m.Lock()
				defer m.Unlock()
				return l.unzipnative(natives)
			}
		case err := <-e:
			return err
		}
	}
}

func (l Libraries) unzipnative(n []string) error {
	e := make(chan error, len(n))
	done := make(chan bool, len(n))
	p := launcher.Minecraft + `/versions/` + l.librarie.ID + `/natives/`
	err := os.MkdirAll(p, 0777)
	if err != nil {
		return err
	}
	go func() {
		for _, v := range n {
			v := v
			go func() {
				err := DeCompress(v, p)
				if err != nil {
					e <- err
				}
				done <- true
			}()
		}
	}()
	i := 0
	for {
		select {
		case <-done:
			i++
			if i == len(n) {
				return nil
			}
		case err := <-e:
			return err
		}
	}
}

func ifallow(l launcher.LibraryX115) bool {
	if l.Rules != nil {
		var allow bool
		for _, r := range l.Rules {
			if r.Action == "disallow" && osbool(r.Os.Name) {
				return false
			}
			if r.Action == "allow" && (r.Os.Name == "" || osbool(r.Os.Name)) {
				allow = true
			}
		}
		return allow
	}
	return true
}

func osbool(os string) bool {
	GOOS := runtime.GOOS
	if GOOS == "darwin" {
		GOOS = "osx"
	}
	return os == GOOS
}

func swichnatives(l launcher.LibraryX115) (path, sha1, url string) {
	Os := runtime.GOOS
	switch Os {
	case "windows":
		path = l.Downloads.Classifiers.NativesWindows.Path
		sha1 = l.Downloads.Classifiers.NativesWindows.Sha1
		url = l.Downloads.Classifiers.NativesWindows.URL
	case "darwin":
		if l.Downloads.Classifiers.NativesOsx.Path != "" {
			path = l.Downloads.Classifiers.NativesOsx.Path
			sha1 = l.Downloads.Classifiers.NativesOsx.Sha1
			url = l.Downloads.Classifiers.NativesOsx.URL
		} else {
			path = l.Downloads.Classifiers.NativesMacos.Path
			sha1 = l.Downloads.Classifiers.NativesMacos.Sha1
			url = l.Downloads.Classifiers.NativesMacos.URL
		}
	case "linux":
		path = l.Downloads.Classifiers.NativesLinux.Path
		sha1 = l.Downloads.Classifiers.NativesLinux.Sha1
		url = l.Downloads.Classifiers.NativesLinux.URL
	default:
		panic("???")
	}
	return
}

func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if !strings.Contains(strings.ToTitle(file.Name), strings.ToTitle("META-INF")) && (strings.HasSuffix(strings.ToTitle(file.Name), strings.ToTitle("dll")) || strings.HasSuffix(strings.ToTitle(file.Name), strings.ToTitle("dylib")) || strings.HasSuffix(strings.ToTitle(file.Name), strings.ToTitle("so"))) {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			filename := dest + file.Name
			if err != nil {
				return err
			}
			w, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer w.Close()
			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
			w.Close()
			rc.Close()
		}
	}
	return nil
}
