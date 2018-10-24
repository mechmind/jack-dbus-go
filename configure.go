package jackdbus

import "github.com/godbus/dbus"

type Configure struct {
	object
	iface string
}

// methods

type ParameterConstraint struct {
	Value dbus.Variant
	Key   string
}

func (config *Configure) GetParameterConstraint(names []string) (isRange bool, isStrict bool, isFakeValue bool, params []ParameterConstraint, err error) {
	err = call(config.Object(), config.iface, "GetParameterConstraint", 0, args{names}, args{&isRange, &isStrict, &isFakeValue, &params})
	if err != nil {
		return false, false, false, nil, err
	}
	return
}

type ParameterInfo struct {
	Flags       byte
	Name        string
	Description string
	Default     string
}

func (config *Configure) GetParameterInfo(names []string) (info ParameterInfo, err error) {
	err = call(config.Object(), config.iface, "GetParameterInfo", 0, args{names}, args{&info})
	return
}

func (config *Configure) GetParameterValue(names []string) (isSet bool, def dbus.Variant, val dbus.Variant, err error) {
	err = call(config.Object(), config.iface, "GetParameterValue", 0, args{names}, args{&isSet, &def, &val})
	return
}

func (config *Configure) SetParameterValue(names []string, value dbus.Variant) error {
	return call(config.Object(), config.iface, "SetParameterValue", 0, args{names, value}, nil)
}

func (config *Configure) ResetParameterValue(names []string) error {
	return call(config.Object(), config.iface, "ResetParameterValue", 0, args{names}, nil)
}

func (config *Configure) GetParametersInfo(parent []string) (info []ParameterInfo, err error) {
	err = call(config.Object(), config.iface, "GetParametersInfo", 0, args{parent}, args{&info})
	return
}

func (config *Configure) ReadContainer(parent []string) (leaf bool, children []string, err error) {
	err = call(config.Object(), config.iface, "ReadContainer", 0, args{parent}, args{&leaf, &children})
	return
}
