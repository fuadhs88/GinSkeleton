package sql_factory

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql 驱动
	"strings"
	//_ "github.com/denisenkom/go-mssqldb" // sqlserver驱动
	// _ "github.com/lib/pq"		//  postgreSql  驱动
	"go.uber.org/zap"
	"goskeleton/app/core/event_manage"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/yml_config"
	"time"
)

var mysqlDriverWrite *sql.DB
var mysqlDriverRead *sql.DB

var sqlServerDriverWrite *sql.DB
var sqlServerDriverRead *sql.DB

var postgreSqlDriverWrite *sql.DB
var postgreSqlDriverRead *sql.DB

// 初始化数据库驱动
func initSqlDriver(sqlType, readOrWrite string) *sql.DB {
	configFac := yml_config.CreateYamlFactory()
	var tmpDriver *sql.DB
	var err error
	if sqlType == "mysql" {
		Host := configFac.GetString("Mysql." + readOrWrite + ".Host")
		Port := configFac.GetString("Mysql." + readOrWrite + ".Port")
		User := configFac.GetString("Mysql." + readOrWrite + ".User")
		Pass := configFac.GetString("Mysql." + readOrWrite + ".Pass")
		Charset := configFac.GetString("Mysql." + readOrWrite + ".Charset")
		DataBase := configFac.GetString("Mysql." + readOrWrite + ".DataBase")
		SetMaxIdleConns := configFac.GetInt("Mysql." + readOrWrite + ".SetMaxIdleConns")
		SetMaxOpenConns := configFac.GetInt("Mysql." + readOrWrite + ".SetMaxOpenConns")
		SetConnMaxLifetime := configFac.GetDuration("Mysql." + readOrWrite + ".SetConnMaxLifetime")
		SqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local&charset=%s", User, Pass, Host, Port, DataBase, Charset)
		switch readOrWrite {
		case "Write", "Read":
			tmpDriver, err = sql.Open("mysql", SqlConnString)
		default:
			variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + readOrWrite)
			return nil
		}
		if err != nil {
			variable.ZapLog.Error(my_errors.ErrorsDbSqlDriverInitFail, zap.Error(err))
			return nil
		}
		tmpDriver.SetMaxIdleConns(SetMaxIdleConns)
		tmpDriver.SetMaxOpenConns(SetMaxOpenConns)
		tmpDriver.SetConnMaxLifetime(SetConnMaxLifetime * time.Second)
		// 将需要销毁的事件统一注册在事件管理器，由程序退出时统一销毁
		event_manage.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Mysql_DB"+readOrWrite, func(args ...interface{}) {
			_ = tmpDriver.Close()
		})
		switch readOrWrite {
		case "Write":
			mysqlDriverWrite = tmpDriver
		case "Read":
			mysqlDriverRead = tmpDriver
		default:
			return nil
		}
		return tmpDriver

	} else if sqlType == "sqlserver" || sqlType == "mssql" {
		Host := configFac.GetString("SqlServer." + readOrWrite + ".Host")
		Port := configFac.GetString("SqlServer." + readOrWrite + ".Port")
		DataBase := configFac.GetString("SqlServer." + readOrWrite + ".DataBase")
		User := configFac.GetString("SqlServer." + readOrWrite + ".User")
		Pass := configFac.GetString("SqlServer." + readOrWrite + ".Pass")
		SetMaxIdleConns := configFac.GetInt("SqlServer." + readOrWrite + ".SetMaxIdleConns")
		SetMaxOpenConns := configFac.GetInt("SqlServer." + readOrWrite + ".SetMaxOpenConns")
		SetConnMaxLifetime := configFac.GetDuration("SqlServer." + readOrWrite + ".SetConnMaxLifetime")
		SqlConnString := fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=disable", Host, Port, DataBase, User, Pass)
		switch readOrWrite {
		case "Write", "Read":
			tmpDriver, err = sql.Open("mssql", SqlConnString)
		default:
			variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + readOrWrite)
			return nil
		}
		if err != nil {
			variable.ZapLog.Error(my_errors.ErrorsDbSqlDriverInitFail + err.Error())
			return nil
		}
		tmpDriver.SetMaxIdleConns(SetMaxIdleConns)
		tmpDriver.SetMaxOpenConns(SetMaxOpenConns)
		tmpDriver.SetConnMaxLifetime(SetConnMaxLifetime * time.Second)
		// 将需要销毁的事件统一注册在事件管理器，由程序退出时统一销毁
		event_manage.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Sqlserver_DB"+readOrWrite, func(args ...interface{}) {
			_ = tmpDriver.Close()
		})
		switch readOrWrite {
		case "Write", "Read":
			sqlServerDriverWrite = tmpDriver
		default:
			return nil
		}
		return tmpDriver
	} else if sqlType == "postgre" {
		Host := configFac.GetString("PostgreSql." + readOrWrite + ".Host")
		Port := configFac.GetString("PostgreSql." + readOrWrite + ".Port")
		DataBase := configFac.GetString("PostgreSql." + readOrWrite + ".DataBase")
		User := configFac.GetString("PostgreSql." + readOrWrite + ".User")
		Pass := configFac.GetString("PostgreSql." + readOrWrite + ".Pass")
		SetMaxIdleConns := configFac.GetInt("PostgreSql." + readOrWrite + ".SetMaxIdleConns")
		SetMaxOpenConns := configFac.GetInt("PostgreSql." + readOrWrite + ".SetMaxOpenConns")
		SetConnMaxLifetime := configFac.GetDuration("PostgreSql." + readOrWrite + ".SetConnMaxLifetime")
		SqlConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", Host, Port, DataBase, User, Pass)
		switch readOrWrite {
		case "Write", "Read":
			tmpDriver, err = sql.Open("postgres", SqlConnString)
		default:
			variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + readOrWrite)
			return nil
		}
		if err != nil {
			variable.ZapLog.Error(my_errors.ErrorsDbSqlDriverInitFail + err.Error())
			return nil
		}
		tmpDriver.SetMaxIdleConns(SetMaxIdleConns)
		tmpDriver.SetMaxOpenConns(SetMaxOpenConns)
		tmpDriver.SetConnMaxLifetime(SetConnMaxLifetime * time.Second)
		// 将需要销毁的事件统一注册在事件管理器，由程序退出时统一销毁
		event_manage.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Postgre_DB"+readOrWrite, func(args ...interface{}) {
			_ = tmpDriver.Close()
		})
		switch readOrWrite {
		case "Write":
			postgreSqlDriverWrite = tmpDriver
		case "Read":
			postgreSqlDriverRead = tmpDriver
		default:
			return nil
		}
		return tmpDriver
	}
	return nil
}

// 从底层驱动中获取一个连接，初始化驱动的过程本质上就是根据参数初始化了一个连接池
func GetOneSqlClient(sqlType, readOrWrite string) *sql.DB {
	if !strings.Contains("mysql,sqlserver,mssql,postgre", sqlType) {
		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + sqlType)
		return nil
	}
	if !strings.Contains("Read,Write", readOrWrite) {
		variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + "，" + readOrWrite)
		return nil
	}
	var maxRetryTimes int
	var reConnectInterval time.Duration
	configFac := yml_config.CreateYamlFactory()

	var dbDriver *sql.DB
	switch sqlType {
	case "mysql":
		if readOrWrite == "Write" {
			if mysqlDriverWrite == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = mysqlDriverWrite
			}
		} else if readOrWrite == "Read" {
			if mysqlDriverRead == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = mysqlDriverRead
			}
		}
		maxRetryTimes = configFac.GetInt("Mysql." + readOrWrite + ".PingFailRetryTimes")
		reConnectInterval = configFac.GetDuration("Mysql." + readOrWrite + ".ReConnectInterval")
	case "sqlserver", "mssql":
		if readOrWrite == "Write" {
			if sqlServerDriverWrite == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = sqlServerDriverWrite
			}
		} else if readOrWrite == "Read" {
			if sqlServerDriverRead == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = sqlServerDriverRead
			}
		}
		maxRetryTimes = configFac.GetInt("SqlServer." + readOrWrite + ".PingFailRetryTimes")
		reConnectInterval = configFac.GetDuration("SqlServer." + readOrWrite + ".ReConnectInterval")
	case "postgre":
		if readOrWrite == "Write" {
			if postgreSqlDriverWrite == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = postgreSqlDriverWrite
			}
		} else if readOrWrite == "Read" {
			if postgreSqlDriverRead == nil {
				dbDriver = initSqlDriver(sqlType, readOrWrite)
			} else {
				dbDriver = postgreSqlDriverRead
			}
		}
		maxRetryTimes = configFac.GetInt("PostgreSql." + readOrWrite + ".PingFailRetryTimes")
		reConnectInterval = configFac.GetDuration("PostgreSql." + readOrWrite + ".ReConnectInterval")
	default:
		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + "，" + sqlType)
		return nil
	}
	if dbDriver == nil {
		return nil
	}
	for i := 1; i <= maxRetryTimes; i++ {
		// ping 失败允许重试
		if err := dbDriver.Ping(); err != nil { //  获取一个连接失败，进行重试
			dbDriver = initSqlDriver(sqlType, readOrWrite)
			time.Sleep(time.Second * reConnectInterval)
			if i == maxRetryTimes {
				variable.ZapLog.Error(sqlType+my_errors.ErrorsDbGetConnFail, zap.Error(err))
				return nil
			}
		} else {
			break
		}
	}
	return dbDriver
}
