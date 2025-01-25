# bookings
Web App from Udemy course from https://github.com/tsawler


Built in Go version 1.22.3
- Uses 	[scs](github.com/alexedwards/scs/v2) v2.8.0
- Uses  [chi router](github.com/go-chi/chi) v1.5.5
- Uses	[nosurf](github.com/justinas/nosurf) v1.1.1
- Uses [GoValidator](https://github.com/asaskevich/govalidator) v0.0.0-20230301143203-a9d515a09cc2
- Uses [GoDotEnv](https://github.com/joho/godotenv) v1.5.1

JavaScript Sources
- Uses [Bootstrap](https://getbootstrap.com/docs/5.3/getting-started/introduction/)
- Uses [Vanilla JS Datepicker](https://mymth.github.io/vanillajs-datepicker/#/)
- Uses [Notie](https://github.com/jaredreich/notie)
- Uses [SweettAlert2](https://sweetalert2.github.io/#download)

Database
- Connects to PostgreSQL backend through driver package [driver.go](https://github.com/cyberjourney20/bookings/blob/main/internal/driver/driver.go) that can be swaped out if another DB is used. 
- Includes database [migrations](https://github.com/cyberjourney20/bookings/tree/main/migrations) generated with [soda](https://gobuffalo.io/documentation/database/soda/)