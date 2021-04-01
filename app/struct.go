package main

import(
	"gorm.io/gorm"
	"time"
)
const(
	CertificateTypeServer	= 1
	CertificateTypeCA	= 2
	CertificateTypeClient	= 3
)
type User struct {
	gorm.Model
	Name	string
	Surname	string
	Org	string
	Mail	string
}
type AccessGroup struct {
	gorm.Model
	Name	string
	Subnet	string
	Mask	string
}
type Certificates struct {
	gorm.Model
	UserId	uint	//ref to User.ID
	CN	string	//Common Name of Cert
	Type	int	//one of CertificateType...
	Private	string	//Private Key
	Public	string	//Public Key
	Revoked	bool	//true, if it is revoked
	ValidTo	time.Time // Valid to flag
}
type UserAccess struct{
	gorm.Model
	UserId		uint	//ref to User.ID
	AccessGroup	uint	//ref to AccessGroup.ID
}
