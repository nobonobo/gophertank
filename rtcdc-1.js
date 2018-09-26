var lpc = new RTCPeerConnection({
  iceServers: [{ url: "stun:stun.l.google.com:19302" }]
});
var rpc = new RTCPeerConnection({
  iceServers: [{ url: "stun:stun.l.google.com:19302" }]
});
var ldc, rdc;
rpc.ondatachannel = ev => {
  console.log("ondatachannel");
  rdc = ev.channel;;
  rdc.onopen = () => {
    console.log("ropen");
  };
  rdc.onclose = () => {
    console.log("rclose");
  };
  rdc.onmessage = ev => {
    console.log("rmsg:" + ev.data);
  };
  rdc.onerror = err => {
    console.log("rerr:" + err);
  };
};
var _ = (lpc.onicecandidate = ice => {
  if (ice.candidate == null) {
    console.log("lpc ice end");
    var sdp = lpc.localDescription;
    rpc.setRemoteDescription(sdp).then(() => {
      rpc.createAnswer().then(sdp => {
        rpc.setLocalDescription(sdp);
      });
    });
  }
});
var _ = (rpc.onicecandidate = ice => {
  if (ice.candidate == null) {
    console.log("rpc ice end");
    var sdp = rpc.localDescription;
    console.log(sdp);
    lpc.setRemoteDescription(sdp);
  }
});
(() => {
  var dc = lpc.createDataChannel("ch");
  ldc = dc;
  ldc.onopen = () => {
    console.log("lopen");
  };
  ldc.onclose = () => {
    console.log("lclose");
  };
  ldc.onmessage = ev => {
    console.log("lmsg:" + ev.data);
  };
  ldc.onerror = err => {
    console.log("lerr:" + err);
  };
  lpc.createOffer().then(sdp => lpc.setLocalDescription(sdp));
})();
