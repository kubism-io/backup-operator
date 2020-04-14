/*
Copyright 2020 Backup Operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kubism-io/backup-operator/pkg/stream"
)

func NewDirDestination(dir string) (stream.Destination, error) {
	return &dirDestination{
		dir: dir,
	}, nil
}

type dirDestination struct {
	dir string
}

func (f *dirDestination) Store(obj stream.Object) error {
	fp := filepath.Join(f.dir, obj.ID)
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, obj.Data)
	return err
}
