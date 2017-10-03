/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"errors"

	"github.com/schollz/org.eclipse.paho.mqtt.golang/packets"
)

func keepalive(c *Client) {
	DEBUG.Println(PNG, "keepalive starting")

	for {
		select {
		case <-c.stop:
			DEBUG.Println(PNG, "keepalive stopped")
			c.workers.Done()
			return
		case <-c.pingTimer.C:
			DEBUG.Println(PNG, "keepalive sending ping")
			ping := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)
			//We don't want to wait behind large messages being sent, the Write call
			//will block until it it able to send the packet.
			ping.Write(c.conn)
			c.pingRespTimer.Reset(c.options.PingTimeout)
		case <-c.pingRespTimer.C:
			CRITICAL.Println(PNG, "pingresp not received, disconnecting")
			c.workers.Done()
			c.internalConnLost(errors.New("pingresp not received, disconnecting"))
			c.pingTimer.Stop()
			return
		}
	}
}
