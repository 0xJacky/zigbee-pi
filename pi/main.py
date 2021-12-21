import json
import os
import time

import serial
import websocket


def on_message(ws, message):
    print(message)


def on_error(ws, error):
    print(error)


def on_close(ws, close_status_code, close_msg):
    print("### closed ###")


def on_open(ws):
    print("Opened connection")
    ws.send(json.dumps({"token": "CF1008B5-7693-42DB-A1E8-87212A23EEA4"}))
    while True:
        ser.flush()
        line = ser.readline()
        line = line.decode('utf-8')
        line = line.rstrip()
        line = line.replace('T&H:', '')
        line = line.split('  ')
        print("T:%s,H:%s" % (line[0], line[1]))
        try:
            ws.send(json.dumps({"temperature": line[0], "humidity": line[1]}))
        finally:
            time.sleep(2)


if __name__ == "__main__":
    device = os.popen('ls /dev/ttyUSB*').readline().rstrip()
    ser = serial.Serial(device, 9600)

    while True:
        try:
            # websocket.enableTrace(True)
            ws = websocket.WebSocketApp("wss://homework.jackyu.cn/zigbee-pi/api/pi",
                                        on_open=on_open,
                                        on_message=on_message,
                                        on_error=on_error,
                                        on_close=on_close)

            ws.run_forever()
        except Exception as e:
            print(e)

        time.sleep(10)

