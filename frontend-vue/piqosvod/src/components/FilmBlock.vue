<template>
  <div class="film-card" v-if="loaded == true">
    <div class="film-card-inner">
      <div class="card-front">
        <img class="poster-image" v-bind:src="getImgURL()" />

        <h4 class="film-title">
          {{ this.filmData.Name }}
        </h4>
      </div>
      <div class="card-back">
        <h1 class="film-title">
          {{ this.filmData.Name }}
        </h1>
        <div class="film-details">
          <h3 class="film-category">{{ this.filmData.Category }}</h3>
          <h3 class="film-description">{{ this.filmData.Description }}</h3>
        </div>
        <div class="watch-film-button" v-on:click="openVideoPlayer">Watch</div>
        <div class="delete-film-button" v-on:click="deleteFilm">Delete</div>
      </div>
    </div>
  </div>
</template>

<script>
//          src="https://upload.wikimedia.org/wikipedia/commons/thumb/5/59/Brainthatwouldntdie_film_poster.jpg/365px-Brainthatwouldntdie_film_poster.jpg"
import axios from "axios";
// import Vue from 'vue'
export default {
  props: ["filmId"],
  inject: ["host"],

  data() {
    return {
      filmData: {
        Name: "Name",
        Desc: "Desc",
        Category: "Category",
        Owner: "Owner",
        imgUrl: "",
      },
      attachmentRecord: [
        {
          id: 1,
          data :null
        },
      ],
      img : null,
      loaded: false,
      token: sessionStorage.token,
    };
  },
  async created() {
    //console.log(this.filmId);
    // GET request using axios with error handling
    await axios
      .get(this.host + "film?id=" + this.filmId)
      .then((response) => (this.filmData = response.data))
      .catch((error) => {
        this.errorMessage = error.message;
        console.error("There was an error!", error);
      });



    this.loaded = true;

    //console.log(this.filmData);
  },
  methods: {
    openVideoPlayer() {
      console.log("Opening video player, filmId = " + this.filmId);
      console.log(window.location.href + "film/" + this.filmId);
      window.location.href = window.location.href + "film/" + this.filmId;
    },
    async deleteFilm() {
      await axios.delete(this.host + "deleteFilm?id=" + this.filmId, {
        headers: {
          token: this.token,
        },
      });
    },

   getImgURL() {
     let url = this.host + "filmPoster?name=" + this.filmId
     return url


      console.log("rekord: ", this.attachmentRecord);
      let record = this.attachmentRecord[0];

      if (record.data == null) {
      //  set(record, "data", null);
      axios
          .get(this.host + "filmPoster?name=" + this.filmId, {
            headers: {
              token: this.token,
            },
          })
          .then((result) => {
              let reader = new FileReader();
              reader.readAsDataURL(result.data); 
              reader.onload = () => {
                  this.img = reader.result;
              }
          //  this.img = result.data
            // record.data = result.data
            // set(record, "data", result.data);
          });
      }
      // console.log(this.attachmentRecord)
      // return this.attachmentRecord[0].data;
      return this.img
    },
  },
};
</script>

<style>
.film-card {
  background: #333;
  color: aliceblue;
  border-radius: 0.375rem;
  overflow: hidden;
  padding-bottom: 0rem;
  background-color: transparent;
  perspective: 1000px;

  min-width: 300px;
  height: 450px;
  margin: 20px;
  position: relative;
}

.film-card-inner {
  position: inherit;
  width: 100%;
  height: 210%;
  text-align: center;
  transition: transform 0.8s;
  transform-style: preserve-3d;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}
.film-card:hover .film-card-inner {
  transform: rotateY(180deg);
}
.card-front,
.card-back {
  background: #333;
  position: absolute;
  width: 100%;
  height: 100%;
  -webkit-backface-visibility: hidden; /* Safari */
  backface-visibility: hidden;
}

.card-back {
  background: #333;
  color: white;
  transform: rotateY(180deg);
  /* display: relative; */
}

.poster-image {
  width: 300px;
  height: 400px;
}
.delete-film-button {
  position: relative;
  margin-top: 10px;
  padding: 5px;
  margin-left: 30%;
  margin-right: 30%;
  color: aliceblue;
  border: 3px solid #555;
  border-radius: 5px;
  cursor: pointer;
  font-weight: bold;
}
.delete-film-button:hover {
  background-color: red;
  color: #333;
  font-weight: bold;
}

.watch-film-button {
  position: relative;
  padding-top: 80%;
  padding: 5px;
  margin-left: 30%;
  margin-right: 30%;
  color: aliceblue;
  border: 3px solid #555;
  border-radius: 5px;
  cursor: pointer;
  font-weight: bold;
}
.watch-film-button:hover {
  background-color: green;
  color: #333;
  font-weight: bold;
}
</style>