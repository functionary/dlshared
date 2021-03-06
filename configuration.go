/**
 * (C) Copyright 2013, Deft Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dlshared

import (
	"os"
	"fmt"
	"github.com/daviddengcn/go-ljson-conf"
)

type Configuration struct {
	Version string
	PidFile string
	Pid int
	Environment string
	Hostname string
	FileName string
	data *ljconf.Conf
}

const confPathKeyPattern = "%s.%s"

func NewConfiguration(fileName string) (*Configuration, error) {

	conf := &Configuration{ FileName : fileName }

	var err error
	if conf.data, err = ljconf.Load(fileName); err != nil { return nil, NewStackError("Unable to load configuration file - error: %v", err) }

	conf.PidFile = conf.data.String("pidFile", "")

	if len(conf.PidFile) == 0 { return nil, NewStackError("Configuration file error - pidFile not set") }

	conf.Environment = conf.data.String("environment", "")

	conf.Version = conf.data.String("version", "")

	conf.Pid = os.Getpid()

	conf.Hostname, err = os.Hostname()
	if err != nil { return nil, NewStackError("Unable to load hostname - error: %v", err) }

	if len(conf.Version) == 0 { return nil, NewStackError("Configuration file error - version not set") }
	if len(conf.Environment) == 0 { return nil, NewStackError("Configuration file error - environment not set") }

	return conf, nil
}

func (self *Configuration) String(key string, def string) string { return self.data.String(key, def) }

func (self *Configuration) StringWithPath(path, key string, def string) string { return self.String(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) Int(key string, def int) int { return self.data.Int(key, def) }

func (self *Configuration) IntWithPath(path, key string, def int) int { return self.Int(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) Bool(key string, def bool) bool { return self.data.Bool(key, def) }

func (self *Configuration) BoolWithPath(path, key string, def bool) bool { return self.Bool(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) Float(key string, def float64) float64 { return self.data.Float(key, def) }

func (self *Configuration) FloatWithPath(path, key string, def float64) float64 { return self.Float(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) StrList(key string, def [] string) []string { return self.data.StringList(key, def) }

func (self *Configuration) IntList(key string, def []int) []int { return self.data.IntList(key, def) }

func (self *Configuration) List(key string, def []interface{}) []interface{} { return self.data.List(key, def) }

func (self *Configuration) ListWithPath(path, key string, def []interface{}) []interface{} { return self.List(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) Interface(key string, def interface{}) interface{} { return self.data.Interface(key, def) }

func (self *Configuration) InterfaceWithPath(path, key string, def interface{}) interface{} { return self.Interface(fmt.Sprintf(confPathKeyPattern, path, key), def) }

func (self *Configuration) EnvironmentIs(expected string) bool { return self.Environment == expected }

