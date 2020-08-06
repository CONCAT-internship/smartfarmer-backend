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
   6. /Records
   7. /RegisterDevice
   8. /RegisterRecipe



* Server domain

  Public DNS(IPv4): `https://asia-northeast1-superfarmers.cloudfunctions.net`



1. /Insert

   데이터베이스에 센서 데이터를 저장합니다.

   | method |  path   |                request                |       response       |
   | :----: | :-----: | :-----------------------------------: | :------------------: |
   | `POST` | /Insert | (JSON) uuid를 포함한 센서 데이터 정보 | (string) 에러 메세지 |

   - Request body 예시
     - uuid: (string) 아두이노 기기의 고유 번호
     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_level: (boolean) 수위
     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
   
   ![sample](https://user-images.githubusercontent.com/29545214/88491674-956fa500-cfdf-11ea-9be0-3cbbc0910614.png)

2. /DailyAverage

   데이터베이스에서 고유번호가 일치하는 기기의 주간 일일 평균 데이터를 반환합니다. (각 데이터는 소숫점 둘째자리에서 반올림)

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
   

![sample](https://user-images.githubusercontent.com/29545214/89208787-feff3d00-d5f7-11ea-8afe-a051e3b1b5e3.png)

3. /RecentStatus

   데이터베이스에서 해당 기기의 최근 상태값을 찾아 반환합니다. (number형 데이터는 소숫점 둘째자리에서 반올림)

   | method |     path      |    request    |      response      |
   | :----: | :-----------: | :-----------: | :----------------: |
   | `GET`  | /RecentStatus | (string) uuid | (JSON) 최근 상태값 |

   - Query string 예시

     - uuid: (string) 아두이노 기기의 고유번호

       `uuid=123e6b776f000c04`

   - Response body 예시

     - temperature: (number) 온도
     - humidity: (number) 습도
     - pH: (number) 배양액 산성 정도
     - ec: (number) 배양액 이온 정도
     - light: (number) 조도
     - liquid_temperature: (number) 수온
     - liquid_level: (boolean) 수위
     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
     - unix_time: (number) 데이터 저장 시각
     - local_time: (timestamp) 데이터 저장 시각 (UTC+0 기준)
     
     ![sample](https://user-images.githubusercontent.com/29545214/89209154-b72ce580-d5f8-11ea-81ce-c95a2ecd450c.png)



4. /Control

   모바일 앱에서 아두이노 기기를 원격 제어하는 데 사용됩니다.

   **desired_status** 컬렉션의 상태를 변화시키고 아두이노 기기에서 데이터베이스의 변화를 감지하면 사용자가 설정한대로 움직이게 됩니다.

   | method |   path   |       request        |       response       |
   | :----: | :------: | :------------------: | :------------------: |
   | `POST` | /Control | (JSON) uuid와 상태값 | (string) 에러 메세지 |

   - Request body 예시

     - uuid: (string) 대상 아두이노 기기의 고유 번호
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

   - Response body 예시

     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
     

   ![sample](https://user-images.githubusercontent.com/29545214/89105841-c3316f80-d45f-11ea-800a-cf970d1b918f.png)

6. /Records

   특정 필드의 최근 기록들을 불러옵니다.

   | method |   path   |              request              |      response      |
   | :----: | :------: | :-------------------------------: | :----------------: |
   | `POST` | /Records | (JSON) uuid와 필드명, 불러올 시간 | (JSON) 조회한 기록 |

   - Request body 예시

     ```json
     {
       "uuid": "756e6b776f000c04",
       "key": "ec",
       "time": 608400
     }
     ```

     - uuid: (string) 조회할 기기의 고유번호
     - key: (string) 조회할 필드명
     - time: (number) 조회할 기간(초 단위)

   

   - Response body 예시

     ```json
     {
       [
       	{
           "ec": 1.8,
           "unix_time":1596695118
         },
       	{
           "ec": 1.9,
           "unix_time": 1596694938
         },
       	{
           "ec": 1.7,
           "unix_time": 1596694935
         },
       	{
           "ec": 1.6,
           "unix_time": 1596694932
         }
       ]
     }
     ```

     - unix_time: (number) 해당 기록의 등록 시간

   

7. /RegisterDevice

   새로운 기기 정보를 등록합니다.

   | method |      path       |               request               |       response       |
   | :----: | :-------------: | :---------------------------------: | :------------------: |
   | `POST` | /RegisterDevice | (JSON) 사용자 정보와 기기 고유 번호 | (string) 에러 메세지 |

   - Request body 예시

     ```json
     {
       "email": "test@example.com",
       "uuid": "756e6b776f000c04"
     }
     ```

     - email: (string) 사용자 이메일
     - uuid: (string) 등록할 기기의 고유 번호

   

8. /RegisterRecipe

   작물 재배 레시피를 등록합니다.

   만약 같은 기기에 기존 레시피가 있었다면 덮어씁니다.

   | method |      path       | request |       response       |
   | :----: | :-------------: | :-----: | :------------------: |
   | `POST` | /RegisterRecipe | (JSON)  | (string) 에러 메세지 |

   - Request body 예시

     ```json
     {
     	"email": "test@example.com",
       "uuid": "756e6b776f000c04",
       "recipe": {
         "crop": "basil",
         "condition": {
           "temperature_min": 25,
           "temperature_max": 30,
           "humidity_min": 50,
           "humidity_max": 60,
           "liquid_temperature": 20,
           "tray_liquid_level": 10,
           "light": 70,
           "light_time": 16,
           "pH_min": 6.0,
           "pH_max": 6.5,
           "ec_min": 1.0,
           "ec_max": 1.5,
           "planting_distance_min_width": 20,
           "planting_distance_min_height": 20,
           "planting_distance_max_width": 25,
           "planting_distance_max_height": 25
         }
       }
    }
    ```
     
     - email: (string) 사용자 이메일
     - uuid: (string) 기기 고유번호
     - crop: (string) 작물 명
     - temperature_min: (number) 온도 최소값
     - temperature_max: (number) 온도 최대값
     - humidity_min: (number) 습도 최소값
     - humidity_max: (number) 습도 최대값
     - liquid_temperature: (number) 양액 온도
     - tray_liquid_level: (number) 양액 물높이
     - light: (number) 조도
     - light_time: (number) 일조 시간
     - pH_min: (number) pH 최소값
     - pH_max: (number) pH 최대값
     - ec_min: (number) 양액 농도 최소값
     - ec_max: (number) 양액 농도 최대값
     - planting_distance_min_width: (number) 재식 거리 최소 가로
     - planting_distance_min_height: (number) 재식 거리 최소 세로
     - planting_distance_max_width: (number) 재식 거리 최대 가로
     - planting_distance_max_height: (number) 재식 거리 최대 세로

