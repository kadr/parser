export CATEGORY_NAME=Автомобильные шины
export LIMIT=100
export PROXY=
export LENTA_COOKIE=Utk_DvcGuid=d3558ef8-445f-ca5e-de21-601065a9f5a2; GrowthBook_user_id=707063d3-687d-049b-985c-440a4e41b528; App_Cache_MPK=mp300-b1de0bac2c257f3257bf5ef2eea4ecbc; App_Cache_CitySlug=moscow; UserSessionId=51d22be6-44d1-98f6-c50f-6df30dec29a5; App_Cache_MissionAddressMode=%7B%22t%22%3A%22pickup%22%2C%22ids%22%3Afalse%2C%22ma%22%3A%7B%22i%22%3A3149%2C%22a%22%3A%220124%22%2C%22t%22%3A%22%D0%A2%D0%9A124%22%2C%22af%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%2C%207-%D1%8F%20%D0%9A%D0%BE%D0%B6%D1%83%D1%85%D0%BE%D0%…430afd8db1ff; App_Cache_City=%7B%22centerLat%22%3A%2255.75322000%22%2C%22centerLng%22%3A%2237.62255200%22%2C%22id%22%3A1%2C%22isDefault%22%3Atrue%2C%22mainDomain%22%3Afalse%2C%22name%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%20%D0%B8%20%D0%9C%D0%9E%22%2C%22slug%22%3A%22moscow%22%7D; agree_with_cookie=true; qrator_ssid=1770357665.395.lTbaFxWSObk0Ttsg-mk6e8l61nldh9q6enriicjcgdidjrjc6; qrator_jsid=1770357667.770.Bq9PxbSMdwY42iBJ-neidbank85v3n442a3brms5j846ap1a1; uwyiert=e1513271-3b7f-a747-146c-a48268775a92

run:
	go build -o parser cmd/main.go && ./parser
