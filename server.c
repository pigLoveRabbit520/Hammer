#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>

#define BUFFER_SIZE 256

int get_bind_addr_socket(unsigned int local_port);
void set0(char *p, size_t size);

int main(int argc, char *argv[])
{
	int sock;
	if ((sock = get_bind_addr_socket(4000)) < 0)
	{
		exit(1);
	}
	printf("udp server started\n");
	int recvLen;
	char recv_buffer[BUFFER_SIZE];
	char send_buffer[BUFFER_SIZE];
	struct sockaddr_in client_addr;
	socklen_t addr_len = sizeof(client_addr);
	set0(send_buffer, BUFFER_SIZE);

	while(1) {
		set0(recv_buffer, BUFFER_SIZE);
		if((recvLen = recvfrom(sock, recv_buffer, BUFFER_SIZE, 0, (struct sockaddr*)&client_addr, &addr_len)) < 0)
		{
			printf("()recvfrom()函数使用失败了.\r\n");
			close(sock);
			exit(1);
		}
		printf("receive data from client: %s\n", recv_buffer);

		sprintf(send_buffer, "IP: %s Port: %d", inet_ntoa(client_addr.sin_addr), client_addr.sin_port);
		size_t send_len = strlen(send_buffer) + 1;
		if(sendto(sock, send_buffer, send_len, 0, (struct sockaddr*)&client_addr, addr_len) != send_len) {
			printf("sendto fail\r\n");
			close(sock);
			exit(0);
		}
	}
	return 0;
}


int get_bind_addr_socket(unsigned int local_port)
{
	struct sockaddr_in addr;
	int sock = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);
	if(sock < 0)
	{
		printf("创建套接字失败了.\n");
		return -1;
	}
	memset(&addr, 0, sizeof(addr));
	addr.sin_family = AF_INET;
	addr.sin_addr.s_addr = htonl(INADDR_ANY);
	addr.sin_port = htons(local_port);

	if(bind(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0)
	{
		printf("bind() 函数使用失败了.\r\n");
		close(sock);
		return -1;
	}
	return sock;
}

void set0(char *p, size_t size)
{
	memset(p, 0, size);
}