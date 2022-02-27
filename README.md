
# piqosvod
### What is piqosvod?
Piqosvod is a small 


# Backend
| Metoda   | Ścieżka URL | Opis                                                    |
|-------------------|----------------------|------------------------------------------------------------------|
| GET               | /videoManifest       | Get video manifest file                               |
| GET               | /videoSegment        | Get single video segment                                   |
| POST              | /bandwidth           | Post update of current video buffer state and last recorded download speed       |
| GET               | /categories          | Get list of film categories                               |
| GET               | /films               | Get list of films                                             |
| GET               | /filmsByCategory     | Get list of films by category                 |
| GET, POST, DELETE | /film                | Manages data of specified film                                 |
| POST              | /filmFile            | Post film file wideo                                           |
| POST              | /initFilmSession     | Initiate new film session                                       |
| POST              | /login               | Post login data                                            |
| GET               | /tokenCheck          | Check if token is still valid                             |
| GET               | /filmQuality         | Get film quality list                          |
| GET               | /filmPoster          | Get film poster                     |

# Frontend
[![Netlify Status](https://api.netlify.com/api/v1/badges/b1661399-dc43-4485-b71a-ceb3c076f04c/deploy-status)](https://app.netlify.com/sites/friendly-haibt-d1928a/deploys)
