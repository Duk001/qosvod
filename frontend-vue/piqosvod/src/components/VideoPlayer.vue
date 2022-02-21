<template>
  <div>
    <video
      ref="videoPlayer"
      class="video-js vjs-default-skin vjs-big-play-centered"
      data-setup="{}"
    ></video>
  </div>
</template>

<script>
import hosts from '../host'
import videojs from "video.js";
import axios from "axios";
export default {
  inject:['host'],
  data() {
    return {
      token: sessionStorage.token,
      player: null,
      bandwidth: 0,
      bufferLength : 0,
      videoOptions: {
        //  autoplay: true,
        controls: true,
        //  preload: true,
        width: "720",
        height: "480",
        sources: [
          {
            src: this.host+"videoManifest?film=",
            type: "application/x-mpegURL",
          },
        ],
        html5: {
          vhs: {
            overrideNative: true,
          },
          nativeAudioTracks: false,
          nativeVideoTracks: false,
        },
      },
      playerLoaded: false,
    };
  },
  async created() {
    
    let str = window.location.pathname;
    console.log("url: " + str);
    let tokens = str.split("/").slice(2, 3);
    let result = tokens.join("/");
    console.log("filmName: " + result);
    let filmName = result;
    this.videoOptions.sources[0].src += filmName;
    this.videoOptions.width = window.innerWidth;
    this.videoOptions.height = window.innerHeight;
    console.log("player settings: " + this.videoOptions);

    var Data = `{"FilmID":"${filmName}"}`;
    var jsonData = JSON.parse(Data);
    console.log("Host: ",this.host)
    const res = await axios.post(
      this.host+"initFilmSession",
      jsonData,
      {
        headers: {
          "Content-Type": "application/json",
          token: this.token,
        },
      }
    );
  },
  mounted() {
    //!!!!!!!!!
    this.player = videojs(
      this.$refs.videoPlayer,
      this.videoOptions,
      function onPlayerReady() {
        //let tmp = this.tech();
        console.log("onPlayerReady", this);

        //setBeforeRequest();
      }
    );
    console.log("paused:", this.player.paused());
    //!!!!!!!!!
    this.playerStateListener();
   // this.vhs = this.player.vhs
  },
  beforeDestroy() {
    if (this.player) {
      //console.log("Żegnam");
      this.player.dispose();
    }
  },
  methods: {
    onPlayerReady() {
      console.log("onPlayerReady", this.player);
      this.playerLoaded = true;
    },

    playerStateListener() {
      const stateListener = async () => {
        // if(this.player.tech({ IWillNotUseThisInPlugins: true }).vhs()){
        //   console.log("TEST: ")
          
        // }
        if (!this.player.paused()) {
          var vhs = this.player.tech({ IWillNotUseThisInPlugins: true }).vhs;
          
          if(this.bandwidth != vhs.bandwidth)
          {
            var bufferEnd = this.player.buffered().end(0)
            var timeStamp = this.player.currentTime()
            this.bufferLength = bufferEnd - timeStamp 
            this.bandwidth = vhs.bandwidth
            var Data = `{"bandwidth":"${this.bandwidth}", "buffer":"${this.bufferLength}"}`;
            var jsonData = JSON.parse(Data);
            //setTimeout(800)
            const res = axios.post(
              this.host+"bandwidth",
              jsonData,
              {
                headers: {
                  "Content-Type": "application/json",
                  token: this.token,
                },
              }
            );

            
          }
        } //else {
          //var vhs = this.player.tech({ IWillNotUseThisInPlugins: true }).vhs;
          //console.log("Paused vhs: ",vhs)
       // }

        return setTimeout(stateListener, 200);
      };
      stateListener();
    },
  },
  watch: {
    // player(newValue, oldValue){
    //   console.log("From Watch: ",oldValue)
    //   console.log("From Watch: ",newValue)
    // }
  },
};
//let counter = 0;

let host = hosts["prod"] //!ustaw prod na produkcję
videojs.options.hls.overrideNative = true;
videojs.Vhs.xhr.beforeRequest = function (options) {
  //debugger
  //console.log("before request intialized, URL: ", options);
  let currentRoute = window.location.pathname;
  let str = currentRoute;
  //console.log("url: " + str);
  let tokens = str.split("/").slice(2, 3);
  let result = tokens.join("/");
  console.log("filmName: " + result);
  let filmName = result;
  if (options.handleManifestRedirects == true) {
     //debugger;
    options.uri = options.uri.replace(
      host+"",
      host+"videoSegment?filmName=" + filmName + "&seg="
    );
    options.headers = {
      Token: sessionStorage.token,
    };
  }
  //console.log(options);
  return options;
};

// playerStateListener();
// function playerStateListener() {
//   const stateListener = () => {
//     if (!player.paused()) {
//       // Player logic -> zmienić ->  wystartować przy zmianie segmentu
//       var vhs = player.tech((IWillNotUseThisInPlugins = true)).vhs;
//       console.log(vhs.bandwidth);

//       var xhr = new XMLHttpRequest();
//       xhr.open("POST", this.host+"bandwidth", true);
//       xhr.setRequestHeader("Content-Type", "application/json");
//       xhr.send(
//         JSON.stringify({
//           bandwidth: vhs.bandwidth,
//           user: "test",
//         })
//       );
//     } else {
//     }
//     return setTimeout(stateListener, 100);
//   };
//   stateListener();
// }
</script>

<style>
@import url("https://vjs.zencdn.net/7.15.4/video-js.css");
.video-js {
  position: absolute !important;
  width: 100% !important;
  height: 100% !important;
}
</style>

//https://codingexplained.com/coding/front-end/vue-js/accessing-vue-instance-outside-declaration