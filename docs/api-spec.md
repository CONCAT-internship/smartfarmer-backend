#### Smart farm API specification

**1. API 명세서 개요**

본 문서는 API 정의와 명세를 포함하며, 개발 목적으로 활용함.



**2. API 목록**

1. /Insert
2. /Get
3. /DailyAverage



**3. API 명세서**

0. Server domain

   Public DNS(IPv4): `https://asia-northeast1-superfarmers.cloudfunctions.net`

   

1. /Insert

   데이터베이스에 센서 데이터를 저장합니다.

   | method |  path   |                request                |
   | :----: | :-----: | :-----------------------------------: |
   | `POST` | /Insert | (JSON) uuid를 포함한 센서 데이터 정보 |

   - Request Body 예시

     - uuid: (string) 아두이노 기기의 고유 번호
     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_flow_rate: (number) 유량
     - liquid_level: (bool) 수위
     - valve: (bool) 밸브 on/off
     - led: (bool) LED on/off
     - fan: (bool) 팬 on/off

     <img width="394" src="https://user-images.githubusercontent.com/29545214/87919744-37106700-cab3-11ea-8b5a-c9f3ada6d1ae.png">

   
<hr>

2. /Get

   데이터베이스에서 고유번호가 일치하는 기기의 최근 일주일 간 내역을 반환합니다.

   | method | path |         request         |                response                |
   | :----: | :--: | :---------------------: | :------------------------------------: |
   | `GET`  | /Get | (string) 기기 고유 번호 | (JSON) 해당 기기의 최근 일주일 간 내역 |

   - Query string 예시

     `uuid=00331447`

   - Response body 예시

     - uuid: (string) 아두이노 기기의 고유 번호
     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_flow_rate: (number) 유량
     - liquid_level: (bool) 수위
     - valve: (bool) 밸브 on/off
     - led: (bool) LED on/off
     - fan: (bool) 팬 on/off
     - unix_time: (number) 해당 레코드의 생성 시간 (unix time)
     - local_time: (timestamp) 해당 레코드의 생성 시간
     
     <img width="732" src="https://user-images.githubusercontent.com/29545214/88047752-0f8dcd00-cb8d-11ea-9705-009f40233864.png">

<hr>

3. /DailyAverage

   데이터베이스에서 고유번호가 일치하는 기기의 일일 데이터 평균을 반환합니다.

   | method | path          | request                            | response                |
   | ------ | ------------- | ---------------------------------- | ----------------------- |
   | `GET`  | /DailyAverage | (string) 고유번호와 날짜 시작 시간 | (json) 일일 데이터 평균 |

   - Query string 예시

     - uuid: 아두이노 기기의 고유번호
     - unixtime: 해당 날짜의 00시 00분 00초의 유닉스 시간

     `uuid=00331447&unixtime=1595289600`

     

   - Response body 예시

     - temperature: 일일 기온 평균
     - humidity: 일일 습도 평균
     - pH: 일일 산성도 평균
     - ec: 일일 이온 농도 평균
     - light: 일일 조도 평균
     - liquid_temperature: 일일 수온 평균
     - liquid_flow_rate: 일일 유랑 평균

     <img width="412" src="https://user-images.githubusercontent.com/29545214/88143835-ddcf4180-cc32-11ea-9a1f-580470e9ffbb.png">

