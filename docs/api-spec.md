#### Smart farm API specification

**1. API 명세서 개요**

본 문서는 API 정의와 명세를 포함하며, 개발 목적으로 활용함.



**2. API 목록**

1. /Insert
2. /Get



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

     `00 33 14 47`

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

