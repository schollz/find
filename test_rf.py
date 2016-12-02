import socket
import time

a = time.time()
for i in range(1):
    # Connect to the server
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('localhost', 5006))

    # Send the data
    message = 'r3=9u6ugjhl4x'
    print ('Sending : "%s"' % message)
    len_sent = s.send(message.encode('utf-8'))

    # Receive a response
    response = s.recv(1024)
    print ('Received: "%s"' % response)

    # Clean up
    s.close()

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('localhost', 5006))

    # Send the data
    message = 'r3='
    print ('Sending : "%s"' % message)
    len_sent = s.send(message.encode('utf-8'))

    # Receive a response
    response = s.recv(1024)
    print ('Received: "%s"' % response)

    # Clean up
    s.close()
print(time.time()-a)
