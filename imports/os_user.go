// this file was generated by gomacro command: import "os/user"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	pkg "os/user"
	. "reflect"
)

func Package_os_user() (map[string]Value, map[string]Type) {
	return map[string]Value{
			"Current":       ValueOf(pkg.Current),
			"Lookup":        ValueOf(pkg.Lookup),
			"LookupGroup":   ValueOf(pkg.LookupGroup),
			"LookupGroupId": ValueOf(pkg.LookupGroupId),
			"LookupId":      ValueOf(pkg.LookupId),
		}, map[string]Type{
			"Group":               TypeOf((*pkg.Group)(nil)).Elem(),
			"UnknownGroupError":   TypeOf((*pkg.UnknownGroupError)(nil)).Elem(),
			"UnknownGroupIdError": TypeOf((*pkg.UnknownGroupIdError)(nil)).Elem(),
			"UnknownUserError":    TypeOf((*pkg.UnknownUserError)(nil)).Elem(),
			"UnknownUserIdError":  TypeOf((*pkg.UnknownUserIdError)(nil)).Elem(),
			"User":                TypeOf((*pkg.User)(nil)).Elem(),
		}
}

func init() {
	binds, types := Package_os_user()
	Binds["os/user"] = binds
	Types["os/user"] = types
}