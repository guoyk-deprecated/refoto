package main

var (
	envPort       = 4000
	envTitle      = "Refoto"
	envDebug      = false
	envMySQLDSN   = ""
	envAdminToken = ""
	envSecret     = ""
)

func setupEnv() (err error) {
	if err = envInt("REFOTO_PORT", &envPort); err != nil {
		return
	}
	if err = envStr("REFOTO_TITLE", &envTitle); err != nil {
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
	return
}
