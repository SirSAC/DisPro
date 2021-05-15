// Made by SirSAC for Network.
package main
//
import (
	. "context"
	"log"
	. "net"
	"runtime"
)
//
func handle_internet(local_connection Conn, target_address string, processor_thread int, pipe_size int, try_count int, keep_alive bool, serial_order bool) {
	if serial_order == true {
		sync_group.Add(processor_thread)
	}
	load_balancer := get_load_balancer(serial_order)
	local_address, _ := ResolveTCPAddr("tcp", load_balancer.interface_address)
	network_dialer.LocalAddr = local_address
	remote_connection, err := network_dialer.DialContext(Background(), "tcp", target_address)
	if err != nil {
		try_again := 0
		for try_again < try_count {
			load_balancer := get_load_balancer(serial_order)
			local_address, _ := ResolveTCPAddr("tcp", load_balancer.interface_address)
			network_dialer.LocalAddr = local_address
			remote_connection, err := network_dialer.DialContext(Background(), "tcp", target_address)
			if err == nil {
				go handle_proxy(local_connection, remote_connection, local_address, target_address, pipe_size, keep_alive)
				if serial_order == true {
					defer sync_group.Wait()
				}
				return
			}
			try_again++
		}
		log.Println(string(COLOR_YELLOW), "[!]", local_address, "<=>", target_address, string(COLOR_RESET))
		local_connection.Write([]byte {5, NETWORK_UNREACHABLE, 0, 1, 0, 0, 0, 0, 0, 0})
		local_connection.Close()
		if serial_order == true {
			defer sync_group.Wait()
			defer sync_group.Done()
		}
		runtime.Goexit()
	}
	go handle_proxy(local_connection, remote_connection, local_address, target_address, pipe_size, keep_alive)
	if serial_order == true {
		defer sync_group.Wait()
	}
}
//
func execute_command(mtu_size string, option_setting bool) {
}
