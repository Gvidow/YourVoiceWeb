export class CustomWebSocket {
    constructor() {
        this.socket = null;
    }

    newConnect(addr, socketOpen, socketClose, socketMessage, socketError) {
        this.socket = new WebSocket(addr);
        this.socket.onopen = socketOpen;
        this.socket.onclose = socketClose;
        this.socket.onmessage = socketMessage;
        this.socket.onerror = socketError;  
    }

    close() {
        console.log("call websocket close");
        this.socket.close();
        this.socket = null;
    }
}