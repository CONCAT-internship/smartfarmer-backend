#### smartfarmer-backend

Back-end module of smartfarmer

![build status](https://github.com/CONCAT-internship/smartfarmer-backend/blob/master/assets/images/badge.svg)



**API specification**

0. Index
   1. /Insert
   2. /DailyAverage
   3. /RecentStatus
   4. /Control
   5. /DesiredStatus



* Server domain

  Public DNS(IPv4): `https://asia-northeast1-superfarmers.cloudfunctions.net`



1. /Insert

   데이터베이스에 센서 데이터를 저장합니다.

   | method |  path   |                request                |       response       |
   | :----: | :-----: | :-----------------------------------: | :------------------: |
   | `POST` | /Insert | (JSON) uuid를 포함한 센서 데이터 정보 | (string) 에러 메세지 |

   - Request Body 예시
     - uuid: (string) 아두이노 기기의 고유 번호
     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_flow_rate: (number) 유량
     - liquid_level: (boolean) 수위
     - valve: (boolean) 밸브 on/off
     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off

   ![sample](https://user-images.githubusercontent.com/29545214/88491674-956fa500-cfdf-11ea-9be0-3cbbc0910614.png)

2. /DailyAverage

   데이터베이스에서 고유번호가 일치하는 기기의 주간 일일 평균 데이터를 반환합니다.

   만약 목요일까지의 데이터만 있고 금요일과 토요일의 데이터는 없다면 금요일과 토요일의 평균 데이터는 비어 있게 됩니다. (빈 구조체)

   | method |     path      |            request             |           response           |
   | :----: | :-----------: | :----------------------------: | :--------------------------: |
   | `GET`  | /DailyAverage | (string) uuid와 주의 시작 시각 | (JSON) 주간 일일 평균 데이터 |

   - Query string 예시
     -	uuid: (string) 아두이노 기기의 고유번호
     -	unixtime: (number) 해당 주의 일요일 0시 0분 0초의 유닉스 시간. (UTC+0 기준)

     `uuid=123e6b776f000c04&unixtime=1595116800`

   - Response body 예시
     - temperature: (number) 일일 기온 평균
     - humidity: (number) 일일 습도 평균
     - pH: (number) 일일 산성도 평균
     - ec: (number) 일일 이온 농도 평균
     - light: (number) 일일 조도 평균
     - liquid_temperature: (number) 일일 수온 평균
     - liquid_flow_rate: (number) 일일 유랑 평균

   ![sample](https://user-images.githubusercontent.com/29545214/88491867-6c501400-cfe1-11ea-95c4-fb808147b413.png)

3. /RecentStatus

   데이터베이스에서 해당 기기의 최근 상태값을 찾아 반환합니다.

   | method |     path      |    request    |      response      |
   | :----: | :-----------: | :-----------: | :----------------: |
   | `GET`  | /RecentStatus | (string) uuid | (JSON) 최근 상태값 |

   - Query string 예시

     - uuid: (string) 아두이노 기기의 고유번호

       `uuid=123e6b776f000c04`

   - Response body 예시
   
     - uuid: (string) 아두이노 기기의 고유 번호
     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_flow_rate: (number) 유량
     - liquid_level: (boolean) 수위
     - valve: (boolean) 밸브 on/off
     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
     
     ![sample](https://user-images.githubusercontent.com/29545214/89105731-b7917900-d45e-11ea-91e3-0ef268125c63.png)



4. /Control

   모바일 앱에서 아두이노 기기를 원격 제어하는 데 사용됩니다.

   **desired_status** 컬렉션의 상태를 변화시키고 아두이노 기기에서 데이터베이스의 변화를 감지하면 사용자가 설정한대로 움직이게 됩니다.

   | method |   path   |       request        |       response       |
   | :----: | :------: | :------------------: | :------------------: |
   | `POST` | /Control | (JSON) uuid와 상태값 | (string) 에러 메세지 |

   - Request body 예시

     - uuid: (string) 대상 아두이노 기기의 고유 번호
     - valve: (boolean) 밸브 on/off
     - led: (boolean) LED를 on/off
     - fan: (boolean) 팬을 on/off

     ![sample](https://user-images.githubusercontent.com/29545214/89105750-ec9dcb80-d45e-11ea-8887-264cbe1d1ef0.png)

   

5. /DesiredStatus

   아두이노 기기에서 사용자의 설정을 리스닝할 때 사용됩니다.

   **desired_status** 컬렉션의 상태를 받아 희망하는 상태값을 파악하고, 기기가 이 결과대로 움직이게 해야합니다.
   
   | method |      path      |    request    |    response    |
| :----: | :------------: | :-----------: | :------------: |
   | `GET`  | /DesiredStatus | (string) uuid | (JSON) 상태 값 |

   - Query string 예시

     `uuid=123e6b776f000c04`

   - Response Body 예시
   
     - valve: (boolean) 밸브 on/off
  - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
   
     ![sample](https://user-images.githubusercontent.com/29545214/89105841-c3316f80-d45f-11ea-800a-cf970d1b918f.png)