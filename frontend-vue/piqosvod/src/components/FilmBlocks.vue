<template>
  <div class="film-blocks">
    <ul class="list">
      <li
        v-for="(data, index) in listOfFilmsJson"
        :key="`data-${index}`"
        class="film-block-elem"
      >
        <film-block :filmId="data"></film-block>
      </li>
    </ul>
  </div>
</template>

<script>
import axios from "axios";
import FilmBlock from "./FilmBlock";
//import FilmBlock from './FilmBlock.vue';
export default {
  
  name: "film-blocks",
  inject:['host'],
  props: ["categories"],
  components: {
    FilmBlock,
    //FilmBlock,
  },
  data() {
    return {
      listOfFilmsJson: null,
    };
  },
  async created() {
    if (this.categories.length == 0) {
      // GET request using axios with error handling
      axios
        .get(this.host+"films")
        .then((response) => (this.listOfFilmsJson = response.data))
        .catch((error) => {
          this.errorMessage = error.message;
          console.error("There was an error!", error);
        });
      //console.log(this.listOfFilmsJson);
    } else {
      // GET request using axios with error handling
      var query = this.categories.join();
      var url = this.host+"filmsByCategory?category=" + query;
      axios
        .get(url)
        .then((response) => (this.listOfFilmsJson = response.data))
        .catch((error) => {
          this.errorMessage = error.message;
          console.error("There was an error!", error);
        });
    }
  },
};
</script>

<style>
.film-blocks {
  margin-top: auto;
  padding: 40px;
}
.film-block-elem {
  margin: auto;
  list-style-type: none;
  padding: 5px;
  margin-top: 10px;
  text-align: center;
  width: 300px;
  height: 450px;
}
.list {
  display: flex;
  flex-flow: row wrap;
  justify-content: space-around;
  padding: 0;
  margin: 0;
  color: #000000;

  list-style: none;
}
</style>