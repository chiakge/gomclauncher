package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xmdhs/gomclauncher/launcher"
)

func Getversionlist(atype string) (*Version, error) {
	var rep *http.Response
	var err error
	var b []byte
	f := auto(atype)
	for i := 0; i < 4; i++ {
		if err := func() error {
			if i == 3 {
				return fmt.Errorf("Getversionlist: %w", err)
			}
			rep, _, err = Aget(source(`https://launchermeta.mojang.com/mc/game/version_manifest.json`, f))
			if rep != nil {
				defer rep.Body.Close()
			}
			if err != nil {
				fmt.Println("获取版本列表失败，重试", fmt.Errorf("Getversionlist: %w", err), source(`https://launchermeta.mojang.com/mc/game/version_manifest.json`, f))
				f = fail(f)
				return nil
			}
			b, err = ioutil.ReadAll(rep.Body)
			if err != nil {
				fmt.Println("获取版本列表失败，重试", fmt.Errorf("Getversionlist: %w", err), source(`https://launchermeta.mojang.com/mc/game/version_manifest.json`, f))
				f = fail(f)
				return nil
			}
			return errors.New("")
		}(); err != nil {
			if err.Error() == "" {
				break
			} else {
				return nil, fmt.Errorf("Getversionlist: %w", err)
			}
		}
	}
	v := Version{}
	err = json.Unmarshal(b, &v)
	v.atype = atype
	if err != nil {
		return nil, fmt.Errorf("Getversionlist: %w", err)
	}
	return &v, nil
}

type Version struct {
	Latest   VersionLatest    `json:"latest"`
	Versions []VersionVersion `json:"versions"`
	atype    string
}

type VersionLatest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type VersionVersion struct {
	ID          string `json:"id"`
	ReleaseTime string `json:"releaseTime"`
	Time        string `json:"time"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func (v Version) Downjson(version string) error {
	f := auto(v.atype)
	for _, vv := range v.Versions {
		if vv.ID == version {
			s := strings.Split(vv.URL, "/")
			path := launcher.Minecraft + `/versions/` + vv.ID + `/` + vv.ID + `.json`
			if ver(path, s[len(s)-2]) {
				return nil
			}
			for i := 0; i < 4; i++ {
				if i == 3 {
					return FileDownLoadFail
				}
				err := get(source(vv.URL, f), path)
				if err != nil {
					fmt.Println("似乎是网络问题，重试", source(vv.URL, f), fmt.Errorf("Downjson: %w", err))
					f = fail(f)
					continue
				}
				if !ver(path, s[len(s)-2]) {
					fmt.Println("文件效验失败，重新下载", source(vv.URL, f))
					f = fail(f)
					continue
				}
				break
			}
			return nil
		}
	}
	return NoSuch
}

var NoSuch = errors.New("no such")
