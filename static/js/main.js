import { CustomWebSocket } from "./custom-web-socket.js";
window.MediaRecorder = OpusMediaRecorder;
const workerOptions = {
    OggOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/OggOpusEncoder.wasm',
    WebMOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/WebMOpusEncoder.wasm'
  };

const btnStart = document.getElementById("record-start");
const btnStop = document.getElementById("record-stop");

let ws = new CustomWebSocket();

navigator.mediaDevices.getUserMedia({audio: true, video: false})
    .then(stream => {
        let options = { mimeType: "audio/wave" };
        const mediaRecorder = new MediaRecorder(stream, options, workerOptions);
        // let voice = [];
        btnStart.addEventListener("click", (ev) => {
            console.log("record start");
            ws.newConnect("ws://localhost:8080", ()=>{console.log("Open new connect");},
                ()=>{console.log("Close connect");}, (mes)=>{console.log(mes);}, (err)=>{alert(err);})
            mediaRecorder.start(100);
        });

        mediaRecorder.addEventListener("dataavailable", (ev) => {
            ws.socket.send(ev.data);
            console.log(`chunk send: size ${ev.data.size}`);
            if (mediaRecorder.state === "inactive") {
                console.log("inactive");
                ws.socket.send("DONE");
                ws.close();
            }
        });

        btnStop.addEventListener("click", (ev) => {
            console.log("record stop")
            mediaRecorder.stop();
            console.log(mediaRecorder.mimeType)
        });
    });
