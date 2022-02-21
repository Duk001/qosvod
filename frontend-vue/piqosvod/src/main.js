import { createApp, h} from 'vue'
import App from './App.vue'
import routes from './routes'
import hosts from './host'



const SimpleRouterApp = {
    provide: {
      //host: "https://piqosvod.azurewebsites.net/"
      host: hosts["prod"]  //! ustaw "prod" na produkcję
    },
    data: () => ({
      currentRoute: window.location.pathname
    }),
  
    computed: {
      ViewComponent () {
        let str = this.currentRoute
        console.log("url: "+str )
        let tokens = str.split("/").slice(0,2)
        let result = tokens.join('/')
        console.log("current route: "+result)
        const matchingPage = routes[result] || '404'
        console.log("Matching page: "+matchingPage)
        //const matchingPage = routes[this.currentRoute] || '404'

        return require(`./pages/${matchingPage}.vue`).default
      }
    },
  
    render () {
        let tmp = h(this.ViewComponent)
        console.log(tmp)
      return h(this.ViewComponent)
    },
  
    created () {
      window.addEventListener('popstate', () => {
        this.currentRoute = window.location.pathname
      })
    }
  }









  // let app=createApp(SimpleRouterApp,{
  //   provide:{
  //     host:"https://piqosvod.azurewebsites.net/"
  //   }
  // }).mount('#app')
createApp(SimpleRouterApp).mount('#app')
//createApp(App).mount('#app')
