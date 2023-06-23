import { CustomWebSocket } from "./custom-web-socket.js";
window.MediaRecorder = OpusMediaRecorder;
const workerOptions = {
    OggOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/OggOpusEncoder.wasm',
    WebMOpusEncoderWasmPath: 'https://cdn.jsdelivr.net/npm/opus-media-recorder@0.8.0/WebMOpusEncoder.wasm'
  };

const btnStart = document.getElementById("record-start");
const btnStop = document.getElementById("record-stop");
const consoleWithQuestionText = document.getElementById("question-text");
const consoleWithAnswerText = document.getElementById("answer-text");
let ws = new CustomWebSocket();
let questionText = "";
let isRecordNow = false;
let isQuestionReceived = false;
let player = document.querySelector('#player');
let voice = [];
let flag = true;

navigator.mediaDevices.getUserMedia({audio: true, video: false})
    .then(stream => {
        let options = { mimeType: "audio/wave" };
        const mediaRecorder = new MediaRecorder(stream, options, workerOptions);
        btnStart.addEventListener("click", (ev) => {
            console.log("record start");
            questionText = "";
            consoleWithQuestionText.innerHTML = "";
            consoleWithAnswerText.innerHTML = "";
            isQuestionReceived = false;
            ws.newConnect("ws://localhost:8080/ws", ()=>{console.log("Open new connect");},
                ()=>{console.log("Close connect");}, (mes) => {
                    // consoleWithQuestionText.innerHTML=mes.data;
                    console.log(mes)
                    console.log(typeof(mes.data));
                    if (isQuestionReceived) {
                        if (typeof mes.data === "string") {
                            consoleWithAnswerText.innerHTML += mes.data;
                        } else {
                            voice.push(mes.data)
                            console.log(typeof(voice))
                            let blob = new Blob(voice, {'type': "audio/wave"});
                            voice = []
                            let audioURL = URL.createObjectURL(blob);
                            player.src = audioURL;
                            player.play()
                            player.onended = () => {ws.socket.send(0)}
                            // player.addEventListener("stop")
                            // (await async () => {
                                // await function playAudio(){
                                //     return new Promise( async res=>{
                                //       player.play()
                                //       player.onended = res
                                //     })
                                //   }()
                            // })();
                            console.log("e");
                            // const newLocal = n = new Audio(audioURL);
                            // n.play();
                        }
                    } else {
                        let res = JSON.parse(mes.data)
                        console.log(res)
                        consoleWithQuestionText.innerHTML = questionText + res.text;
                        if (res.fix) {
                            questionText += res.text;
                        }
                        if (res.finish) {
                            consoleWithQuestionText.insertAdjacentHTML("beforeend", `<h4>Ожидайте ответ на ваш вопрос</h4><p>${questionText}</p>`);
                            isQuestionReceived = true;
                        }
                    }
                    // consoleWithAnswerText.innerHTML += mes.data;
                }, (err)=>{alert(err);})
            mediaRecorder.start(100);
        });

        mediaRecorder.addEventListener("start", (ev) => {
            console.log("FFFFF")
            voice = []
        })

        mediaRecorder.addEventListener("dataavailable", (ev) => {
            // voice.push(ev.data)
            ws.socket.send(ev.data);
            console.log(`chunk send: size ${ev.data.size}`);
            if (mediaRecorder.state === "inactive") {
                console.log("inactive");
                ws.socket.send("DONE");
            }
        });

        btnStop.addEventListener("click", (ev) => {
            console.log("record stop")
            mediaRecorder.stop();
            console.log(mediaRecorder.mimeType)
        });
    });
