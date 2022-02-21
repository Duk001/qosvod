(function(e){function t(t){for(var o,a,c=t[0],s=t[1],l=t[2],f=0,d=[];f<c.length;f++)a=c[f],Object.prototype.hasOwnProperty.call(i,a)&&i[a]&&d.push(i[a][0]),i[a]=0;for(o in s)Object.prototype.hasOwnProperty.call(s,o)&&(e[o]=s[o]);u&&u(t);while(d.length)d.shift()();return r.push.apply(r,l||[]),n()}function n(){for(var e,t=0;t<r.length;t++){for(var n=r[t],o=!0,c=1;c<n.length;c++){var s=n[c];0!==i[s]&&(o=!1)}o&&(r.splice(t--,1),e=a(a.s=n[0]))}return e}var o={},i={app:0},r=[];function a(t){if(o[t])return o[t].exports;var n=o[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,a),n.l=!0,n.exports}a.m=e,a.c=o,a.d=function(e,t,n){a.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},a.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},a.t=function(e,t){if(1&t&&(e=a(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(a.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var o in e)a.d(n,o,function(t){return e[t]}.bind(null,o));return n},a.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return a.d(t,"a",t),t},a.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},a.p="/";var c=window["webpackJsonp"]=window["webpackJsonp"]||[],s=c.push.bind(c);c.push=t,c=c.slice();for(var l=0;l<c.length;l++)t(c[l]);var u=s;r.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},"026f":function(e,t,n){},"02e2":function(e,t,n){"use strict";var o=n("7a23"),i={key:0,class:"film-card"},r={class:"film-card-inner"},a={class:"card-front"},c=["src"],s={class:"film-title"},l={class:"card-back"},u={class:"film-title"},f={class:"film-details"},d={class:"film-category"},b={class:"film-description"};function m(e,t,n,m,p,h){return 1==p.loaded?(Object(o["j"])(),Object(o["e"])("div",i,[Object(o["f"])("div",r,[Object(o["f"])("div",a,[Object(o["f"])("img",{class:"poster-image",src:h.getImgURL()},null,8,c),Object(o["f"])("h4",s,Object(o["o"])(this.filmData.Name),1)]),Object(o["f"])("div",l,[Object(o["f"])("h1",u,Object(o["o"])(this.filmData.Name),1),Object(o["f"])("div",f,[Object(o["f"])("h3",d,Object(o["o"])(this.filmData.Category),1),Object(o["f"])("h3",b,Object(o["o"])(this.filmData.Description),1)]),Object(o["f"])("div",{class:"watch-film-button",onClick:t[0]||(t[0]=function(){return h.openVideoPlayer&&h.openVideoPlayer.apply(h,arguments)})},"Watch"),Object(o["f"])("div",{class:"delete-film-button",onClick:t[1]||(t[1]=function(){return h.deleteFilm&&h.deleteFilm.apply(h,arguments)})},"Delete")])])])):Object(o["d"])("",!0)}var p=n("1da1"),h=(n("96cf"),n("bc3a")),g=n.n(h),j={props:["filmId"],inject:["host"],data:function(){return{filmData:{Name:"Name",Desc:"Desc",Category:"Category",Owner:"Owner",imgUrl:""},attachmentRecord:[{id:1,data:null}],img:null,loaded:!1,token:sessionStorage.token}},created:function(){var e=this;return Object(p["a"])(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,g.a.get(e.host+"film?id="+e.filmId).then((function(t){return e.filmData=t.data})).catch((function(t){e.errorMessage=t.message,console.error("There was an error!",t)}));case 2:e.loaded=!0;case 3:case"end":return t.stop()}}),t)})))()},methods:{openVideoPlayer:function(){console.log("Opening video player, filmId = "+this.filmId),console.log(window.location.href+"film/"+this.filmId),window.location.href=window.location.href+"film/"+this.filmId},deleteFilm:function(){var e=this;return Object(p["a"])(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,g.a.delete(e.host+"deleteFilm?id="+e.filmId,{headers:{token:e.token}});case 2:case"end":return t.stop()}}),t)})))()},getImgURL:function(){var e=this.host+"filmPoster?name="+this.filmId;return e}}},O=(n("b2e5"),n("6b0d")),v=n.n(O);const y=v()(j,[["render",m]]);t["a"]=y},"06f9":function(e,t,n){"use strict";var o=n("7a23"),i={class:"film-blocks"},r={class:"list"};function a(e,t,n,a,c,s){var l=Object(o["n"])("film-block");return Object(o["j"])(),Object(o["e"])("div",i,[Object(o["f"])("ul",r,[(Object(o["j"])(!0),Object(o["e"])(o["a"],null,Object(o["m"])(c.listOfFilmsJson,(function(e,t){return Object(o["j"])(),Object(o["e"])("li",{key:"data-".concat(t),class:"film-block-elem"},[Object(o["g"])(l,{filmId:e},null,8,["filmId"])])})),128))])])}var c=n("1da1"),s=(n("96cf"),n("a15b"),n("bc3a")),l=n.n(s),u=n("02e2"),f={name:"film-blocks",inject:["host"],props:["categories"],components:{FilmBlock:u["a"]},data:function(){return{listOfFilmsJson:null}},created:function(){var e=this;return Object(c["a"])(regeneratorRuntime.mark((function t(){var n,o;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:0==e.categories.length?l.a.get(e.host+"films").then((function(t){return e.listOfFilmsJson=t.data})).catch((function(t){e.errorMessage=t.message,console.error("There was an error!",t)})):(n=e.categories.join(),o=e.host+"filmsByCategory?category="+n,l.a.get(o).then((function(t){return e.listOfFilmsJson=t.data})).catch((function(t){e.errorMessage=t.message,console.error("There was an error!",t)})));case 1:case"end":return t.stop()}}),t)})))()}},d=(n("2d30"),n("6b0d")),b=n.n(d);const m=b()(f,[["render",a]]);t["a"]=m},"0f2d":function(e,t,n){},"0fc0":function(e,t,n){"use strict";n("0f2d")},1:function(e,t){},"1b11":function(e,t,n){"use strict";n("dc79")},"1e05":function(e,t,n){},"2d30":function(e,t,n){"use strict";n("6b7e")},"3bb0":function(e,t,n){"use strict";n("5726")},"40bb":function(e,t,n){"use strict";n("da30")},"43f2":function(e,t,n){"use strict";n.r(t);var o=n("7a23"),i={class:"body"};function r(e,t,n,r,a,c){var s=Object(o["n"])("video-player");return Object(o["j"])(),Object(o["e"])("div",i,[Object(o["g"])(s)])}var a={ref:"videoPlayer",class:"video-js vjs-default-skin vjs-big-play-centered","data-setup":"{}"};function c(e,t,n,i,r,c){return Object(o["j"])(),Object(o["e"])("div",null,[Object(o["f"])("video",a,null,512)])}var s=n("1da1"),l=(n("96cf"),n("fb6a"),n("ac1f"),n("1276"),n("a15b"),n("99af"),n("5319"),n("9764")),u=n("f0e2"),f=n("bc3a"),d=n.n(f),b={inject:["host"],data:function(){return{token:sessionStorage.token,player:null,bandwidth:0,bufferLength:0,videoOptions:{controls:!0,width:"720",height:"480",sources:[{src:this.host+"videoManifest?film=",type:"application/x-mpegURL"}],html5:{vhs:{overrideNative:!0},nativeAudioTracks:!1,nativeVideoTracks:!1}},playerLoaded:!1}},created:function(){var e=this;return Object(s["a"])(regeneratorRuntime.mark((function t(){var n,o,i,r,a,c;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return n=window.location.pathname,console.log("url: "+n),o=n.split("/").slice(2,3),i=o.join("/"),console.log("filmName: "+i),r=i,e.videoOptions.sources[0].src+=r,e.videoOptions.width=window.innerWidth,e.videoOptions.height=window.innerHeight,console.log("player settings: "+e.videoOptions),a='{"FilmID":"'.concat(r,'"}'),c=JSON.parse(a),console.log("Host: ",e.host),t.next=15,d.a.post(e.host+"initFilmSession",c,{headers:{"Content-Type":"application/json",token:e.token}});case 15:t.sent;case 16:case"end":return t.stop()}}),t)})))()},mounted:function(){
//!!!!!!!!!
this.player=Object(u["a"])(this.$refs.videoPlayer,this.videoOptions,(function(){console.log("onPlayerReady",this)})),console.log("paused:",this.player.paused()),//!!!!!!!!!
this.playerStateListener()},beforeDestroy:function(){this.player&&this.player.dispose()},methods:{onPlayerReady:function(){console.log("onPlayerReady",this.player),this.playerLoaded=!0},playerStateListener:function(){var e=this,t=function(){var n=Object(s["a"])(regeneratorRuntime.mark((function n(){var o,i,r,a,c;return regeneratorRuntime.wrap((function(n){while(1)switch(n.prev=n.next){case 0:return e.player.paused()||(o=e.player.tech({IWillNotUseThisInPlugins:!0}).vhs,e.bandwidth!=o.bandwidth&&(i=e.player.buffered().end(0),r=e.player.currentTime(),e.bufferLength=i-r,e.bandwidth=o.bandwidth,a='{"bandwidth":"'.concat(e.bandwidth,'", "buffer":"').concat(e.bufferLength,'"}'),c=JSON.parse(a),d.a.post(e.host+"bandwidth",c,{headers:{"Content-Type":"application/json",token:e.token}}))),n.abrupt("return",setTimeout(t,200));case 2:case"end":return n.stop()}}),n)})));return function(){return n.apply(this,arguments)}}();t()}},watch:{}},m=l["a"]["prod"];//!ustaw prod na produkcję
u["a"].options.hls.overrideNative=!0,u["a"].Vhs.xhr.beforeRequest=function(e){var t=window.location.pathname,n=t,o=n.split("/").slice(2,3),i=o.join("/");console.log("filmName: "+i);var r=i;return 1==e.handleManifestRedirects&&(e.uri=e.uri.replace(m+"",m+"videoSegment?filmName="+r+"&seg="),e.headers={Token:sessionStorage.token}),e};n("7ab7");var p=n("6b0d"),h=n.n(p);const g=h()(b,[["render",c]]);var j=g,O={inject:["host"],components:{VideoPlayer:j}};n("a27b");const v=h()(O,[["render",r]]);t["default"]=v},4444:function(e,t,n){"use strict";n("7a7a")},"56d7":function(e,t,n){"use strict";n.r(t);n("e260"),n("e6cf"),n("cca6"),n("a79d"),n("fb6a"),n("ac1f"),n("1276"),n("a15b");var o=n("7a23");n("02e2");var i=n("06f9"),r=n("d000");r["a"],i["a"],n("ffa9"),n("6b0d");var a={"/":"Home","/film":"VideoPlayerPage"},c=n("9764"),s={provide:{host:c["a"]["prod"]},data:function(){return{currentRoute:window.location.pathname}},computed:{ViewComponent:function(){var e=this.currentRoute;console.log("url: "+e);var t=e.split("/").slice(0,2),o=t.join("/");console.log("current route: "+o);var i=a[o]||"404";return console.log("Matching page: "+i),n("aa59")("./".concat(i,".vue")).default}},render:function(){var e=Object(o["h"])(this.ViewComponent);return console.log(e),Object(o["h"])(this.ViewComponent)},created:function(){var e=this;window.addEventListener("popstate",(function(){e.currentRoute=window.location.pathname}))}};Object(o["b"])(s).mount("#app")},5726:function(e,t,n){},"6b7e":function(e,t,n){},"718e":function(e,t,n){"use strict";n("993f")},"7a7a":function(e,t,n){},"7ab7":function(e,t,n){"use strict";n("adad")},9764:function(e,t,n){"use strict";t["a"]={prod:"https://piqosvod.azurewebsites.net/",dev:"http://127.0.0.1:10000/"}},"993f":function(e,t,n){},a27b:function(e,t,n){"use strict";n("026f")},aa59:function(e,t,n){var o={"./404.vue":"ee5d","./Home.vue":"bc13","./VideoPlayerPage.vue":"43f2"};function i(e){var t=r(e);return n(t)}function r(e){if(!n.o(o,e)){var t=new Error("Cannot find module '"+e+"'");throw t.code="MODULE_NOT_FOUND",t}return o[e]}i.keys=function(){return Object.keys(o)},i.resolve=r,e.exports=i,i.id="aa59"},ac35:function(e,t,n){},adad:function(e,t,n){},b2e5:function(e,t,n){"use strict";n("d4b7")},bc13:function(e,t,n){"use strict";n.r(t);var o=n("7a23"),i={class:"body"};function r(e,t,n,r,a,c){var s=Object(o["n"])("login"),l=Object(o["n"])("nav-bar"),u=Object(o["n"])("category-select"),f=Object(o["n"])("film-blocks");return Object(o["j"])(),Object(o["e"])("div",i,[null==this.token?(Object(o["j"])(),Object(o["c"])(s,{key:0,onLoginToken:t[0]||(t[0]=function(e){return c.assignToken(e)})})):Object(o["d"])("",!0),null!=this.token?(Object(o["j"])(),Object(o["c"])(l,{key:1,onChangeCategoriesVisibility:t[1]||(t[1]=function(e){return c.changeCategoriesVisibility(e)})})):Object(o["d"])("",!0),null!=this.token&&a.categorySelectVisible?(Object(o["j"])(),Object(o["c"])(u,{key:2,onListOfCategories:t[2]||(t[2]=function(e){return c.changeCategories(e)}),onChangeCategoriesVisibility:t[3]||(t[3]=function(e){return c.changeCategoriesVisibility(e)})})):Object(o["d"])("",!0),null!=this.token?(Object(o["j"])(),Object(o["c"])(f,{categories:this.filmCategories,key:a.filmBlocksKey},null,8,["categories"])):Object(o["d"])("",!0)])}var a=n("1da1"),c=(n("96cf"),n("bc3a")),s=n.n(c),l=(n("02e2"),n("06f9")),u=n("d000"),f={class:"login-container"},d={class:"login-form-container"},b={class:"login-form"},m={class:"login-form-list"},p=Object(o["f"])("label",{for:"Login"},"Login:",-1),h=Object(o["f"])("label",{for:"Password"},"Password:",-1);function g(e,t,n,i,r,a){return Object(o["j"])(),Object(o["e"])("div",f,[Object(o["f"])("div",d,[Object(o["f"])("form",b,[Object(o["f"])("ul",m,[Object(o["f"])("li",null,[p,Object(o["r"])(Object(o["f"])("input",{type:"text","onUpdate:modelValue":t[0]||(t[0]=function(e){return r.login=e}),id:"Login",name:"Login"},null,512),[[o["q"],r.login]])]),Object(o["f"])("li",null,[h,Object(o["r"])(Object(o["f"])("input",{type:"password","onUpdate:modelValue":t[1]||(t[1]=function(e){return r.password=e}),id:"Password",name:"Password"},null,512),[[o["q"],r.password]]),Object(o["f"])("div",{class:"login-button",onClick:t[2]||(t[2]=function(e){return a.getToken()})},"Login")])])])])])}n("99af");var j={inject:["host"],data:function(){return{login:null,password:null,token:null}},methods:{getToken:function(){var e=this;return Object(a["a"])(regeneratorRuntime.mark((function t(){var n,o,i,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return o='{"login":"'.concat(e.login,'",\n            "password":"').concat(e.password,'" }'),i=JSON.parse(o),t.next=4,s.a.post(e.host+"login",i,{headers:{"Content-Type":"application/json"}});case 4:r=t.sent,n=r.data,n&&"200"==r.status&&(e.token=n,sessionStorage.token=e.token,sessionStorage.username=e.login,e.$emit("loginToken",e.token),console.log("token: ",e.token));case 7:case"end":return t.stop()}}),t)})))()}}},O=(n("4444"),n("6b0d")),v=n.n(O);const y=v()(j,[["render",g]]);var w=y,k=(n("caad"),n("2532"),{class:"categories-select-body"}),C={class:"list"},F=["onClick"],x={class:"cat-buttons"};function V(e,t,n,i,r,a){return Object(o["j"])(),Object(o["e"])("div",k,[Object(o["f"])("ul",C,[(Object(o["j"])(!0),Object(o["e"])(o["a"],null,Object(o["m"])(r.listOfCategories,(function(e,t){return Object(o["j"])(),Object(o["e"])("li",{key:"data-".concat(t),class:Object(o["i"])([r.selectedCategories.includes(e)?"film-category-elem-selected":"film-category-elem"]),onClick:function(t){return a.toggleCategory(e)}},Object(o["o"])(e),11,F)})),128))]),Object(o["f"])("div",x,[Object(o["f"])("div",{class:"reset-button",onClick:t[0]||(t[0]=function(e){return a.resetCategory()})},"Reset"),Object(o["f"])("div",{class:"close-button",onClick:t[1]||(t[1]=function(e){return a.closeCategory()})},"Close")])])}n("a434");var R={inject:["host"],data:function(){return{listOfCategories:[],selectedCategories:[]}},created:function(){var e=this;return Object(a["a"])(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,s.a.get(e.host+"categories").then((function(t){return e.listOfCategories=t.data})).catch((function(t){e.errorMessage=t.message,console.error("There was an error!",t)}));case 2:case"end":return t.stop()}}),t)})))()},methods:{toggleCategory:function(e){if(this.selectedCategories.includes(e)){var t=this.selectedCategories.indexOf(e);-1!==t&&this.selectedCategories.splice(t,1)}else this.selectedCategories.push(e);this.$emit("listOfCategories",this.selectedCategories)},resetCategory:function(){this.selectedCategories=[],this.$emit("listOfCategories",this.selectedCategories)},closeCategory:function(){this.$emit("changeCategoriesVisibility",!1)}}};n("3bb0");const S=v()(R,[["render",V]]);var P=S,T={name:"App",inject:["host"],components:{NavBar:u["a"],FilmBlocks:l["a"],CategorySelect:P,Login:w},data:function(){return{token:null,login:null,categorySelectVisible:!1,filmCategories:[],filmBlocksKey:0}},mounted:function(){var e=this;return Object(a["a"])(regeneratorRuntime.mark((function t(){var n,o;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(console.log("Host ",e.host),!sessionStorage.token){t.next=8;break}return n={token:sessionStorage.token},t.next=5,s.a.get(e.host+"tokenCheck",{headers:n});case 5:o=t.sent,console.log(o),o.data==sessionStorage.token?e.token=sessionStorage.token:e.token=null;case 8:case"end":return t.stop()}}),t)})))()},methods:{assignToken:function(e){this.token=e},changeCategories:function(e){this.filmCategories=e,this.filmBlocksKey+=1},changeCategoriesVisibility:function(e){this.categorySelectVisible=e}}};n("1b11");const N=v()(T,[["render",r]]);t["default"]=N},d000:function(e,t,n){"use strict";var o=n("7a23"),i=function(e){return Object(o["l"])("data-v-4df82eb6"),e=e(),Object(o["k"])(),e},r={class:"navbar"},a=i((function(){return Object(o["f"])("li",null,[Object(o["f"])("a",{href:"/"},"Home")],-1)})),c=i((function(){return Object(o["f"])("div",{class:"bar1"},null,-1)})),s=i((function(){return Object(o["f"])("div",{class:"bar2"},null,-1)})),l=i((function(){return Object(o["f"])("div",{class:"bar3"},null,-1)})),u=[c,s,l];function f(e,t,n,i,c,s){var l=Object(o["n"])("slide-menu"),f=Object(o["n"])("new-film-form");return Object(o["j"])(),Object(o["e"])(o["a"],null,[Object(o["f"])("div",r,[Object(o["f"])("ul",null,[a,Object(o["f"])("li",null,[Object(o["f"])("div",{class:"cat-button",onClick:t[0]||(t[0]=function(e){return s.openCategorySelect()})},"Categories")]),Object(o["f"])("li",null,[Object(o["f"])("div",{class:"add-button",onClick:t[1]||(t[1]=function(){return s.toogleNewFilmFrom&&s.toogleNewFilmFrom.apply(s,arguments)})},"Add")]),Object(o["f"])("li",{class:"settings",onClick:t[2]||(t[2]=function(){return s.toggleOvelray&&s.toggleOvelray.apply(s,arguments)})},u)])]),Object(o["g"])(l,{isOverlay:c.overlay},null,8,["isOverlay"]),Object(o["g"])(f,{isVisible:c.newFilmForm,onChangeFormVisibility:t[3]||(t[3]=function(e){return s.changeFV(e)})},null,8,["isVisible"])],64)}n("a4d3"),n("e01a");var d={key:1,class:"form-container"},b=Object(o["f"])("div",{class:"card-title"},"Upload Form",-1),m={class:"film-form"},p={class:"form-list"},h=Object(o["f"])("label",{for:"Title"},"Title:",-1),g=Object(o["f"])("label",{for:"Category"},"Category:",-1),j=Object(o["f"])("label",{for:"Description"},"Description:",-1),O=Object(o["f"])("label",{for:"v-model-quality-checkboxes"},"Quality:",-1),v={id:"v-model-quality-checkboxes",class:"quality-checkboxes"},y=["id","value"],w=["for"],k=Object(o["f"])("label",{for:"file"},"Filename:",-1),C=Object(o["f"])("label",{for:"poster"},"Poster:",-1);function F(e,t,n,i,r,a){return Object(o["j"])(),Object(o["e"])(o["a"],null,[n.isVisible?(Object(o["j"])(),Object(o["e"])("div",{key:0,class:"overlay",onClick:t[0]||(t[0]=function(){return a.closeForm&&a.closeForm.apply(a,arguments)})})):Object(o["d"])("",!0),n.isVisible?(Object(o["j"])(),Object(o["e"])("div",d,[b,Object(o["f"])("form",m,[Object(o["f"])("ul",p,[Object(o["f"])("li",null,[h,Object(o["r"])(Object(o["f"])("input",{type:"text","onUpdate:modelValue":t[1]||(t[1]=function(e){return r.title=e}),id:"Title",name:"Title"},null,512),[[o["q"],r.title]])]),Object(o["f"])("li",null,[g,Object(o["r"])(Object(o["f"])("input",{type:"text","onUpdate:modelValue":t[2]||(t[2]=function(e){return r.category=e}),id:"Category",name:"Category"},null,512),[[o["q"],r.category]])]),Object(o["f"])("li",null,[j,Object(o["r"])(Object(o["f"])("textarea",{class:"desc-testarea",type:"text","onUpdate:modelValue":t[3]||(t[3]=function(e){return r.description=e}),id:"Description",name:"Description"},null,512),[[o["q"],r.description]])]),Object(o["f"])("li",null,[O,Object(o["f"])("div",v,[(Object(o["j"])(!0),Object(o["e"])(o["a"],null,Object(o["m"])(r.qualityList,(function(e){return Object(o["j"])(),Object(o["e"])("li",{key:e.quality},[Object(o["r"])(Object(o["f"])("input",{type:"checkbox",id:e.quality,value:e.quality,"onUpdate:modelValue":t[4]||(t[4]=function(e){return r.videoQuality=e})},null,8,y),[[o["p"],r.videoQuality]]),Object(o["f"])("label",{for:e.quality},Object(o["o"])(e.quality),9,w)])})),128))])]),Object(o["f"])("li",null,[k,Object(o["f"])("input",{class:"file-picker",type:"file",name:"file",id:"file",onChange:t[5]||(t[5]=function(){return a.onFileChangeFilm&&a.onFileChangeFilm.apply(a,arguments)})},null,32)]),Object(o["f"])("li",null,[C,Object(o["f"])("input",{class:"poster-picker",type:"file",name:"poster",id:"poster",onChange:t[6]||(t[6]=function(){return a.onFileChangePoster&&a.onFileChangePoster.apply(a,arguments)})},null,32)])])]),Object(o["f"])("button",{class:"submit-button",onClick:t[7]||(t[7]=function(){return a.validateForm&&a.validateForm.apply(a,arguments)})},"Submit")])):Object(o["d"])("",!0)],64)}var x=n("1da1"),V=(n("96cf"),n("99af"),n("bc3a")),R=n.n(V),S={name:"NewFilmForm",inject:["host"],emits:["changeFormVisibility"],props:["isVisible"],data:function(){return{title:null,category:null,description:null,owner:null,filmFile:null,posterFile:null,visible:this.isVisible,qualityList:[{quality:"1920:1080-4000k"},{quality:"1280:720-2500k"},{quality:"1280:720-1700k"},{quality:"960:540-900k"},{quality:"640:360-550k"},{quality:"480:270-350k"}],videoQuality:[]}},methods:{onFileChangeFilm:function(e){var t=e.target.files||e.dataTransfer.files;t.length&&(this.filmFile=t[0])},onFileChangePoster:function(e){var t=e.target.files||e.dataTransfer.files;t.length&&(this.posterFile=t[0])},validateForm:function(){this.sendForm()},sendForm:function(){var e=this;return Object(x["a"])(regeneratorRuntime.mark((function t(){var n,o,i,r,a,c,s,l;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(n='{"title":"'.concat(e.title,'",\n            "owner":"').concat(e.owner,'",\n            "category":"').concat(e.category,'",\n            "description":"').concat(e.description,'",\n            "quality" : "').concat(e.videoQuality,'" }'),o=JSON.parse(n),console.log(o),null==e.filmFile){t.next=21;break}return t.next=6,R.a.post(e.host+"film",o,{headers:{"Content-Type":"application/json"}});case 6:return i=t.sent,r=i.data,a=new FormData,a.append("file",e.filmFile),t.next=12,R.a.post(e.host+"filmFile?name="+r,a,{headers:{"Content-Type":"multipart/form-data"}});case 12:return c=t.sent,console.log("Response(film Data): ",c.status),s=new FormData,s.append("file",e.posterFile),t.next=18,R.a.post(e.host+"filmPoster?name="+r,s,{headers:{"Content-Type":"multipart/form-data"}});case 18:l=t.sent,console.log("Response(film poster): ",l.status),200==i.status&&200==c.status&&(e.visible=!1,e.$emit("changeFormVisibility",e.visible));case 21:case"end":return t.stop()}}),t)})))()},closeForm:function(){this.visible=!1,this.$emit("changeFormVisibility",this.visible)}}},P=(n("40bb"),n("6b0d")),T=n.n(P);const N=T()(S,[["render",F]]);var q=N,L={key:0,class:"container"},D={class:"menu"},I={class:"menu-list"},U=Object(o["f"])("li",null,"About",-1);function M(e,t,n,i,r,a){return n.isOverlay?(Object(o["j"])(),Object(o["e"])("div",L,[Object(o["f"])("div",{class:"overlay",onClick:t[1]||(t[1]=function(t){return e.NOisOverlay=!1})},[Object(o["f"])("div",D,[Object(o["f"])("ul",I,[Object(o["f"])("li",null,Object(o["o"])(r.username),1),U,Object(o["f"])("li",null,[Object(o["f"])("div",{class:"logout-btn",onClick:t[0]||(t[0]=function(e){return a.logout()})},"Logout")])])])])])):Object(o["d"])("",!0)}var _={inject:["host"],props:["isOverlay"],data:function(){return{overlay:!0,newFilmForm:!1,username:null}},mounted:function(){this.username=sessionStorage.username},methods:{logout:function(){var e=this;return Object(x["a"])(regeneratorRuntime.mark((function t(){var n,o;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return n={token:sessionStorage.token},t.next=3,R.a.get(e.host+"logout",{headers:n});case 3:o=t.sent,console.log(o),sessionStorage.username=null,sessionStorage.token=null,window.location.reload(!0);case 8:case"end":return t.stop()}}),t)})))()}}};n("0fc0");const J=T()(_,[["render",M]]);var B=J,E={name:"NavBar",inject:["host"],components:{SlideMenu:B,NewFilmForm:q},data:function(){return{overlay:!1,newFilmForm:!1}},methods:{toggleOvelray:function(){console.log("Overlay"),this.overlay=!this.overlay,1==this.overlay?document.documentElement.style.overflow="hidden":document.documentElement.style.overflow="auto"},toogleNewFilmFrom:function(){this.newFilmForm=!this.newFilmForm,1==this.newFilmForm?document.documentElement.style.overflow="hidden":document.documentElement.style.overflow="auto"},changeFV:function(e){this.newFilmForm=e,1==this.newFilmForm?document.documentElement.style.overflow="hidden":document.documentElement.style.overflow="auto"},openCategorySelect:function(){this.$emit("changeCategoriesVisibility",!0)}}};n("e3ae");const $=T()(E,[["render",f],["__scopeId","data-v-4df82eb6"]]);t["a"]=$},d4b7:function(e,t,n){},da30:function(e,t,n){},dc79:function(e,t,n){},e3ae:function(e,t,n){"use strict";n("1e05")},ee5d:function(e,t,n){"use strict";n.r(t);var o=n("7a23"),i=function(e){return Object(o["l"])("data-v-6f5c7277"),e=e(),Object(o["k"])(),e},r=i((function(){return Object(o["f"])("h1",{class:"msg"},"Page not found. :(",-1)})),a=[r];function c(e,t,n,i,r,c){return Object(o["j"])(),Object(o["e"])("div",null,a)}var s=n("d000"),l={inject:["host"],components:{NavBar:s["a"]}},u=(n("718e"),n("6b0d")),f=n.n(u);const d=f()(l,[["render",c],["__scopeId","data-v-6f5c7277"]]);t["default"]=d},ffa9:function(e,t,n){"use strict";n("ac35")}});
//# sourceMappingURL=app.0c09276e.js.map