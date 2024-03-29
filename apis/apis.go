/*
This package contains shared config of all internal api calls to be used by packages in this folder.
The Viper config object should be built from a YAML structured like this:

  apis:
	internal:
		auth.module.url.base: "https://.../auth"
		media.module.url.base: "https://.../media/api"
		user.module.url.base: "https://.../user/api"
		...

Hence, for each subfolder (package), there should be a correspondong entry
  apis.internal.[name].module.url.base
set to the API endpoint, without a trailing slash (/)

To use the API with the config, call Init() with a valid config object.

If any API calls is used without Init(), the call panic instead
*/
package apis

import "github.com/spf13/viper"

var v *viper.Viper

func V() *viper.Viper {
	if v == nil {
		panic("internal API call without calling Init()")
	}
	return v
}

// Initialize apis with a config object
func Init(config *viper.Viper) {
	v = config
}
