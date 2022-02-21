<template>
  <div class="overlay" v-if="isVisible" v-on:click="closeForm"></div>
  <div class="form-container" v-if="isVisible">
    <div class="card-title">Upload Form</div>
    <form class="film-form">
      <ul class="form-list">
        <li>
          <label for="Title">Title:</label>
          <input type="text" v-model="title" id="Title" name="Title" />
        </li>
        <li>
          <label for="Category">Category:</label>
          <input type="text" v-model="category" id="Category" name="Category" />
        </li>
        <li>
          <label for="Description">Description:</label>
          <textarea
            class="desc-testarea"
            type="text"
            v-model="description"
            id="Description"
            name="Description"
          />
        </li>
        <li>
          <label for="v-model-quality-checkboxes">Quality:</label>
          <div id="v-model-quality-checkboxes" class="quality-checkboxes">
            <li v-for="item in qualityList" :key="item.quality">
              <input
                type="checkbox"
                :id="item.quality"
                :value="item.quality"
                v-model="videoQuality"
              />
              <label :for="item.quality">
                {{ item.quality }}
              </label>
            </li>
          </div>
        </li>
        <li>
          <label for="file">Filename:</label>
          <input
            class="file-picker"
            type="file"
            name="file"
            id="file"
            v-on:change="onFileChangeFilm"
          />
        </li>
        <li>
          <label for="poster">Poster:</label>
          <input
            class="poster-picker"
            type="file"
            name="poster"
            id="poster"
            v-on:change="onFileChangePoster"
          />
        </li>
      </ul>
    </form>
    <button class="submit-button" v-on:click="validateForm">Submit</button>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "NewFilmForm",
  inject: ["host"],
  emits: ["changeFormVisibility"],
  props: ["isVisible"],
  data() {
    return {
      title: null,
      category: null,
      description: null,
      owner: null,
      filmFile: null,
      posterFile: null,
      // posterImageFile: null,
      visible: this.isVisible,
      qualityList: [
        { quality: "1920:1080-4000k" },
        { quality: "1280:720-2500k" },
        { quality: "1280:720-1700k" },
        { quality: "960:540-900k" },
        { quality: "640:360-550k" },
        { quality: "480:270-350k" },
      ],
      videoQuality: [],
    };
  },
  methods: {
    onFileChangeFilm(e) {
      var files = e.target.files || e.dataTransfer.files;
      if (!files.length) {
        return;
      }
      this.filmFile = files[0];
    },
    onFileChangePoster(e) {
      var files = e.target.files || e.dataTransfer.files;
      if (!files.length) {
        return;
      }
      this.posterFile = files[0];
    },
    validateForm() {
      //TODO Validation <---

      this.sendForm();
    },
    async sendForm() {
      var Data = `{"title":"${this.title}",
            "owner":"${this.owner}",
            "category":"${this.category}",
            "description":"${this.description}",
            "quality" : "${this.videoQuality}" }`;
      var jsonData = JSON.parse(Data);
      console.log(jsonData);
      //return
      if (this.filmFile != null) {
        //debugger
        //? File data
        const res = await axios.post(this.host + "film", jsonData, {
          headers: {
            "Content-Type": "application/json",
          },
        });
        //console.log("Response(film Data): ",res.status)
        //console.log(res)

        var filmId = res.data;
        
        
        //? File
        const formData = new FormData();
        formData.append("file", this.filmFile);
        const res2 = await axios.post(
          this.host + "filmFile?name=" + filmId,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        console.log("Response(film Data): ", res2.status);
        
        const formData_poster = new FormData();
        formData_poster.append("file", this.posterFile);
        const res3 = await axios.post(
          this.host + "filmPoster?name=" + filmId,
          formData_poster,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        console.log("Response(film poster): ", res3.status);
        if (res.status == 200 && res2.status == 200) {
          this.visible = false;
          this.$emit("changeFormVisibility", this.visible);
        }
        //res.data.data;
      }
    },
    closeForm() {
      this.visible = false;
      this.$emit("changeFormVisibility", this.visible);
      // this.isVisible = false
      //this.isVisible = false
    },
  },
};
</script>

<style>
.form-container {
  position: fixed;
  top: 50%;
  left: 50%;
  z-index: 2;
  /* bring your own prefixes */
  transform: translate(-50%, -50%);
  background-color: #333;
  color: aliceblue;
  border: 3px solid #444;
  border-radius: 5px;
}
.film-form {
  padding: 10px 0px 10px 0px;
}
.form-list {
  list-style-type: none;
  text-align: left;
}
.form-list label {
  display: block;
}
.card-title {
  margin: 10px;
  font-weight: bold;
  font-size: x-large;
}
.desc-testarea {
  overflow: auto;
  resize: none;
  background-color: #f8f8f8;
}
.submit-button {
  margin-bottom: 10px;
  background-color: #555;
  color: aliceblue;
  border-radius: 5px;
  font-size: large;
  border: 3px solid #444;
  cursor: pointer;
}
.submit-button:hover {
  background: #666;
  border: 3px solid #555;
}
.overlay {
  position: absolute;
  display: block;
  width: 100%;
  height: 100%;
  background-color: rgba(51, 51, 51, 0.705);
  z-index: 1;
  cursor: pointer;
}
</style>