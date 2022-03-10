| API      | Description | path
| ----------- | ----------- | ----------- | 
| alert create      | creates an alert       | /api/v1/create-alert
| alert details   | returns the details of an alert        | /api/alert-details
| alert enable | Returns a string with a success message if alert is enabled/disabled successfully,or an error message if the request fails to enable/disable the alert.| /api/alert-enable
| alert remove | Returns a string with a success message if the alert is deleted successfully,or an error message if the request fails to delete the alert. | /api/alert-remove
| alerts list | retruns a json object containing a list of alerts| /api/alerts-list
|influxdb write v1||
|influx db write v2||
|create metric group|Returns a string with success message if the metric group was successfully created,or an error message if the request fails to create the metric group.|/api/metric-group-create
|metric group details|Returns a json string containing the details of the requested metric group.| /api/metric-group-details
|metric group remove| Returns a string with a success message if metric group is deleted successfully,or an error message if the request fails to delete the metric group.|/api/metric-group-remove
|metric groups list|Returns a json object containing a list of metric groups.|/api/metric-group-list
|notification group create|Returns a string with success message if the notification group was successfully created,or an error message if the request fails to create the notification group.|/api/notification-group-create
|notification group details|Returns a json string containing the details of the requested notification group.|/api/notification-group-details
|notification group remove| Returns a string with a success message if the notification group is deleted successfully,or an error message if the request fails to delete the notification group.|/api/notification-group-remove
|notification groups list|Returns json containing a list of notification groups.|/api/notification-groups-list
|open tsdb put||/api/put
|suspension create|Returns a string with success message if the suspension was successfully created,or an error message if the request fails to create the suspension.|/api/suspension-create
|suspension details|Returns a json string containing the details of the requested suspension.|/api/suspension-details
|suspension enable|Returns a string with a success message if suspension is enabled/disabled successfully,or an error message if the request fails to enable/disable the suspension.|/api/suspension-enable
|suspension remove|Returns a string with a success message if the suspension is deleted successfully,or an error message if the request fails to delete the suspension.|/api/suspension-remove
|suspension lists|Returns json containing a list of suspensions|/api/suspensions-list