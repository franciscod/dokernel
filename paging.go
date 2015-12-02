package main

import (
	"reflect"

	"github.com/franciscod/godo"
)

func getList(client *godo.Client, t interface{}, dropletID int) (interface{}, error) {
	vt := reflect.ValueOf(t)
	ts := reflect.SliceOf(vt.Type())
	list := reflect.MakeSlice(ts, 0, 0)

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		var _elems interface{}
		var resp *godo.Response
		var err error

		switch t.(type) {
		case godo.Action:
			_elems, resp, err = client.Actions.List(opt)
		case godo.Droplet:
			_elems, resp, err = client.Droplets.List(opt)
		case godo.Kernel:
			_elems, resp, err = client.Droplets.Kernels(dropletID, opt)
		}

		if err != nil {
			return nil, err
		}

		ve := reflect.ValueOf(_elems)

		for i := 0; i < ve.Len(); i++ {
			list = reflect.Append(list, ve.Index(i))
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}

	return list.Interface(), nil
}

func dropletList(client *godo.Client) ([]godo.Droplet, error) {
	list, err := getList(client, godo.Droplet{}, 0)
	return list.([]godo.Droplet), err
}

func actionList(client *godo.Client) ([]godo.Action, error) {
	list, err := getList(client, godo.Action{}, 0)
	return list.([]godo.Action), err
}

func kernelList(client *godo.Client, dropletID int) ([]godo.Kernel, error) {
	list, err := getList(client, godo.Kernel{}, dropletID)
	return list.([]godo.Kernel), err
}
