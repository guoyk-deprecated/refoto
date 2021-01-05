package main

var (
	envPort               = 4000
	envTitle              = "Refoto"
	envContact            = ""
	envDebug              = false
	envMySQLDSN           = ""
	envAdminToken         = ""
	envSecret             = ""
	envOSSBucket          = ""
	envOSSEndpoint        = ""
	envOSSAccessKeyID     = ""
	envOSSAccessKeySecret = ""
	envOSSPublicEndpoint  = ""
)

func setupEnv() (err error) {
	if err = envInt("REFOTO_PORT", &envPort); err != nil {
		return
	}
	if err = envStr("REFOTO_TITLE", &envTitle); err != nil {
		return
	}
	if err = envStr("REFOTO_CONTACT", &envContact); err != nil {
		return
	}
	if err = envBool("REFOTO_DEBUG", &envDebug); err != nil {
		return
	}
	if err = envStr("REFOTO_MYSQL_DSN", &envMySQLDSN); err != nil {
		return
	}
	if err = envStr("REFOTO_SECRET", &envSecret); err != nil {
		return
	}
	if err = envStr("REFOTO_ADMIN_TOKEN", &envAdminToken); err != nil {
		return
	}
	if err = envStr("REFOTO_OSS_BUCKET", &envOSSBucket); err != nil {
		return
	}
	if err = envStr("REFOTO_OSS_ENDPOINT", &envOSSEndpoint); err != nil {
		return
	}
	if err = envStr("REFOTO_OSS_AK_ID", &envOSSAccessKeyID); err != nil {
		return
	}
	if err = envStr("REFOTO_OSS_AK_SECRET", &envOSSAccessKeySecret); err != nil {
		return
	}
	if err = envStr("REFOTO_OSS_PUBLIC_ENDPOINT", &envOSSPublicEndpoint); err != nil {
		return
	}
	return
}
