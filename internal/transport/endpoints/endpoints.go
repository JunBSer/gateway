package endpoints

import "github.com/JunBSer/gateway/internal/metadata"

func SetupEndpoints(cfgEnd *metadata.EndpointConfig) {
	cfgEnd.AddEndpoint("GET", "/v1/hotels", "Hotels.ListHotels", metadata.AuthUser)
	cfgEnd.AddEndpoint("GET", "/v1/hotels/{hotel_id}/rooms", "Rooms.ListRooms", metadata.AuthUser)
	cfgEnd.AddEndpoint("GET", "/v1/hotels/{hotel_id}/rooms/{id}", "Rooms.GetRoom", metadata.AuthUser)

	cfgEnd.AddEndpoint("POST", "/v1/auth/login", "Auth.Login", metadata.AuthNone)
	cfgEnd.AddEndpoint("POST", "/v1/auth/register", "Auth.Register", metadata.AuthNone)
	cfgEnd.AddEndpoint("POST", "/v1/auth/refresh", "Auth.RefreshToken", metadata.AuthNone)

	cfgEnd.AddEndpoint("GET", "/v1/hotels/{id}", "Hotels.GetHotel", metadata.AuthUser)
	cfgEnd.AddEndpoint("GET", "/v1/hotels/search", "Hotels.SearchHotels", metadata.AuthUser)
	cfgEnd.AddEndpoint("GET", "/v1/hotels/{hotel_id}/availability", "Hotels.GetAvailability", metadata.AuthUser)

	cfgEnd.AddEndpoint("GET", "/docs/swagger.json", "Docs.GetSwagger", metadata.AuthNone)

	cfgEnd.AddEndpoint("POST", "/v1/auth/logout", "Auth.Logout", metadata.AuthUser)
	cfgEnd.AddEndpoint("PUT", "/v1/users/me/password", "Auth.ChangePassword", metadata.AuthUser)
	cfgEnd.AddEndpoint("DELETE", "/v1/users/me", "Auth.DeleteAccount", metadata.AuthUser)
	cfgEnd.AddEndpoint("PUT", "/v1/users/me", "Auth.UpdateProfile", metadata.AuthUser)

	cfgEnd.AddEndpoint("POST", "/v1/admin/users", "Auth.CreateUser", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("GET", "/v1/admin/users/{user_id}", "Auth.GetUser", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("GET", "/v1/admin/users", "Auth.ListUsers", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("PUT", "/v1/admin/users/{user_id}", "Auth.UpdateUser", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("DELETE", "/v1/admin/users/{user_id}", "Auth.DeleteUser", metadata.AuthAdmin)

	cfgEnd.AddEndpoint("POST", "/v1/hotels", "Hotels.CreateHotel", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("PUT", "/v1/hotels/{id}", "Hotels.UpdateHotel", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("DELETE", "/v1/hotels/{id}", "Hotels.DeleteHotel", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("POST", "/v1/hotels/{hotel_id}/rooms", "Rooms.CreateRoom", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("PUT", "/v1/hotels/{hotel_id}/rooms/{id}", "Rooms.UpdateRoom", metadata.AuthAdmin)
	cfgEnd.AddEndpoint("DELETE", "/v1/hotels/{hotel_id}/rooms/{id}", "Rooms.DeleteRoom", metadata.AuthAdmin)

	cfgEnd.AddEndpoint("POST", "/v1/bookings", "BookingService.CreateBooking", metadata.AuthUser)
	cfgEnd.AddEndpoint("POST", "/v1/getbooking", "BookingService.GetBooking", metadata.AuthUser)
	cfgEnd.AddEndpoint("DELETE", "/v1/cancel/{booking_id}", "BookingService.CancelBooking", metadata.AuthUser)
	cfgEnd.AddEndpoint("GET", "/v1/admin/bookings", "BookingService.ListBookings", metadata.AuthAdmin)
}
