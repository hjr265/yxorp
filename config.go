// Copyright 2014 The Yxorp Authors. All rights reserved.

package main

import (
	"io"
	"os"

	"github.com/hjr265/go-zrsc/zrsc"
	"github.com/pelletier/go-toml"
)

var cfg *toml.TomlTree

func init() {
	_, err := os.Stat("config.tml")
	if os.IsNotExist(err) {
		f2, err := zrsc.Open("config-sample.tml")
		catch(err)

		f, err := os.Create("config.tml")
		_, err = io.Copy(f, f2)
		catch(err)

		err = f2.Close()
		catch(err)
		err = f.Close()
		catch(err)
	}

	cfg, err = toml.LoadFile("config.tml")
	catch(err)
}
