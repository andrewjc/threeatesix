package common


type GeneralProtectionFault struct {

}

func (GeneralProtectionFault) Error() string {
	return "General Protection Fault"
}

