from datetime import datetime
from dateutil import tz
 
# consider date in string format
 
# convert datetime string into date,month,day and
# hours:minutes:and seconds format using strptime
def convert_time(time):
  from_zone = tz.tzutc()
  to_zone = tz.tzlocal()

  utc = datetime.strptime(time, "%Y/%m/%d %H:%M:%S")

  utc = utc.replace(tzinfo=from_zone)

  central = utc.astimezone(to_zone)

  central = datetime.strftime(central, "%Y/%m/%d %H:%M:%S")

  print(central)




