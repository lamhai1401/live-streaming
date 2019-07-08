// get DOM elements
var dataChannelLog = document.getElementById('data-channel'),
  iceConnectionLog = document.getElementById('ice-connection-state'),
  iceGatheringLog = document.getElementById('ice-gathering-state'),
  signalingLog = document.getElementById('signaling-state');

// peer connection
var pc = null;
// data channel
var dc = null, dcInterval = null;

async function request_answer(idRoom, idPeer) {
  let answer_1 = await pc.createAnswer();
  console.log("Answer ***: ")
  console.log(answer_1);
  console.log(pc);
  await pc.setLocalDescription(answer_1);
  await fetch('/answer', {
    headers: {
      'Content-Type': 'application/json; charset=utf-8'
    },
    method: 'POST',
    body: JSON.stringify({
        "idRoom": idRoom,
        "idPeer": idPeer,
        "session": {
            "sdp": answer_1.sdp,
            "type": answer_1.type
        }
    })
  })

}

async function join() {
  try {
    test_id = null
    document.getElementById('join').style.display = 'none';
    let idRoom = document.getElementById("idRoom").value;

    let { offer , idPeer} = await fetch('/viewer', {
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "idRoom": idRoom
      }),
      method: 'POST'
    }).then(e => e.json())    
    console.log(offer);    
    var config = {
      sdpSemantics: 'unified-plan',
    };

    pc = new RTCPeerConnection(config);

    // connect audio / video
    pc.addEventListener('track', function (evt) {
      console.log("Track event: ", evt)
      if (evt.track.kind == 'video')
        document.getElementById('video').srcObject = evt.streams[0];
      else
        document.getElementById('audio').srcObject = evt.streams[0];
    });

    console.log('Preparing to set remote');

    await pc.setRemoteDescription({
      sdp: offer.sdp,
      type: offer.type
    });

    await request_answer(idRoom, idPeer);


    console.log("ICE promise ***: ")

    // wait for ICE gathering to complete
    await new Promise(function (resolve) {
      if (pc.iceGatheringState === 'complete') {
        resolve();
      } else {
        function checkState() {
          if (pc.iceGatheringState === 'complete') {
            pc.removeEventListener('icegatheringstatechange', checkState);
            resolve();
            console.log("icegatheringstatechange complete")
          }
        }
        pc.addEventListener('icegatheringstatechange', checkState);
      }
    });

    console.log("finished answer");
  } catch (error) {
    alert(error)
  }

}