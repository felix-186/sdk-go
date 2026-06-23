package driver

import (
	"errors"
	"io"
	"strings"

	"github.com/felix-186/logger"
)

type ErrorType int

const (
	UNKONWN                     ErrorType = 1
	TIMEOUT                     ErrorType = 2
	CONNECTION_FAIELD           ErrorType = 3
	CONNECTION_CLOSED           ErrorType = 4
	CONNECTION_EOF              ErrorType = 5
	MODBUS_TRANSACTION          ErrorType = 6
	MODBUS_ILLEGAL_DATA_ADDRESS ErrorType = 7
)

func TcpClientErrSuggest(err error) (ErrorType, error) {
	if strings.Contains(err.Error(), "timeout") {
		return TIMEOUT, logger.NewErrorFocusNotice("检查网络是否存在延迟；检查服务端设备资源(CPU、内存等)占用是否过高，可尝试降低采集频率", err)
	} else if strings.Contains(err.Error(), "An established connection was aborted by the software in your host machine") ||
		strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") ||
		strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") ||
		strings.Contains(err.Error(), "connection refused") {
		return CONNECTION_FAIELD, logger.NewErrorFocusNotice("检查服务端设备是否开机；检查网络端口是否连通；检查防火墙端口是否开放", err)
	} else if strings.Contains(err.Error(), "broken pipe") ||
		strings.Contains(err.Error(), "use of closed network connection") ||
		strings.Contains(err.Error(), "no route to host") ||
		strings.Contains(err.Error(), "connection reset by peer") {
		return CONNECTION_CLOSED, logger.NewErrorFocusNotice("检查服务端设备是否超过最大连接数限制；检查连接是否空闲超时被断开", err)
	} else if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return CONNECTION_EOF, logger.NewErrorFocusNotice("检查服务端设备是否主动关闭了连接；检查连接是否被服务端重置", err)
	}
	return UNKONWN, err
}

func ModbusErrSuggest(err error) (ErrorType, error) {
	if strings.Contains(err.Error(), "modbus: response transaction id") && strings.Contains(err.Error(), "does not match request") {
		return MODBUS_TRANSACTION, logger.NewErrorFocusNotice("检查是否存在多个连接同时读写导致事务ID不匹配；建议使用连接池或加锁控制并发", err)
	} else if strings.Contains(err.Error(), "illegal data address") {
		return MODBUS_ILLEGAL_DATA_ADDRESS, logger.NewErrorFocusNotice("检查站号和寄存器地址是否配置正确；确认寄存器地址范围是否在设备支持范围内", err)
	} else {
		return TcpClientErrSuggest(err)
	}
}
