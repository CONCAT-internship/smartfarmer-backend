#### smartfarmer-backend

Back-end module of smartfarmer

![build status](https://github.com/CONCAT-internship/smartfarmer-backend/blob/master/assets/images/badge.svg)



**API specification**

0. Index
   1. /Insert
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

   데이터 검수 과정 중 이상값이 감지되면 **abnormal** 컬렉션과 **desired_status** 컬렉션의 값을 업데이트합니다.

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

     ```json
     {
       "uuid": "756e6b776f000c04",
       "temperature": 27.4,
       "humidity": 66.1,
       "pH": 6.3,
       "ec": 0.9,
       "light": 70,
       "liquid_temperature": 25,
       "liquid_level": false,
       "led": true,
       "fan": true
     }
     ```

     

     데이터에 이상값이 감지되면 **abnormal** 컬렉션엔 다음과 같은 정보가 기록됩니다.

     ```json
     {
       "uuid": "756e6b776f000c04",
       "errors": [
         5011,
         5020,
         ...
       ]
     }
     ```

     - uuid: (string) 아두이노 기기의 고유 번호
     - errors: (Array&lt;number&gt;) 에러 코드

     다음은 `errors` 필드의 에러 코드들에 대한 명세입니다.

     ```bash
     7000: CODE_DATA_EMPTY (기기에서 값을 보내지 못했음. 기기 고장)
     
     4000: CODE_PH_MALFUNC (기기에서 보낸 pH 값이 측정 범위를 벗어남. pH 센서 고장)
     4001: CODE_EC_MALFUNC (기기에서 보낸 ec 값이 측정 범위를 벗어남. ec 센서 고장)
     4002: CODE_LIGHT_MALFUNC (기기에서 보낸 조도 값이 측정 범위를 벗어남. 조도 센서 고장)
     
     5000: CODE_PH_IMPROPER_HIGH (pH 값이 레시피에서 정한 범위 초과. pH 펌프를 닫아 pH를 낮춤)
     5001: CODE_PH_IMPROPER_LOW (pH 값이 레시피에서 정한 범위 미만. pH 펌프를 열어 PH를 높임)
     5010: CODE_EC_IMPROPER_HIGH (ec 값이 레시피에서 정한 범위 초과. ec 펌프를 열어 ec를 높임)
     5011: CODE_EC_IMPROPER_LOW (ec 값이 레시피에서 정한 범위 미만)
     5020: CODE_TEMPERATURE_IMPROPER_HIGH (온도가 레시피에서 정한 범위 초과. fan을 가동시켜 온도를 낮춤)
     5021: CODE_TEMPERATURE_IMPROPER_LOW (온도가 레시피에서 정한 범위 미만. fan을 중지시켜 온도를 높임)
     5030: CODE_HUMIDITY_IMPROPER_HIGH (습도가 레시피에서 정한 범위 초과)
     5031: CODE_HUMIDITY_IMPROPER_LOW (습도가 레시피에서 정한 범위 미만)
     ```

     

     또한 **desired_status** 컬렉션의 document id가 uuid와 일치하는 문서를 다음과 같이 업데이트합니다.

     ```json
     {
       "led": false,
       "fan": true,
       "ph_pump": 1,
       "ec_pump": 0
     }
     ```

     이에 대한 내용은 `DesiredStatus` API를 호출하여 확인할 수 있습니다.

     `4. /DesiredStatus` 를 참고해주세요.

     


2. /RecentStatus

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
     
     ```json
     {
       "temperature": 27.4,
       "humidity": 66.1,
       "pH": 6.3,
       "ec": 1.3,
       "light": 70,
       "liquid_temperature": 25,
       "liquid_level": false,
       "led": true,
       "fan": true,
       "unix_time": 1596694932,
       "local_time": "2020-08-09 10:46:13.613475 +0900 KST m=+0.000062068"
     }
     ```
     
     

3. /Control

   모바일 앱에서 아두이노 기기를 원격 제어하는 데 사용됩니다.

   **desired_status** 컬렉션의 상태를 변화시키고 아두이노 기기에서 데이터베이스의 변화를 감지하면 사용자가 설정한대로 움직이게 됩니다.

   | method |   path   |       request        |       response       |
   | :----: | :------: | :------------------: | :------------------: |
   | `POST` | /Control | (JSON) uuid와 상태값 | (string) 에러 메세지 |

   - Request body 예시

     - uuid: (string) 대상 아두이노 기기의 고유 번호
     - led: (boolean) LED를 on/off
     - fan: (boolean) 팬을 on/off
     
     ```json
     {
       "uuid": "756e6b776f000c04",
       "status": {
         "led": false,
         "fan": true
       }
     }
     ```
     
     

4. /DesiredStatus

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
     - pH_pump: (number) pH 펌프 조절 정도 (1: pH 증가, 0: 유지, -1: pH 감소)
     - ec_pump: (number) ec 펌프 조절 정도 (1: ec 증가, 0: 유지)
     
     ```json
     {
       "led": false,
       "fan": true,
       "pH_pump": 1,
       "ec_pump": 0
     }
     ```
     
     

5. /Records

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
           "ec": 1.5,
           "local_time": "2020-08-09 10:46:13.613475 +0900 KST m=+0.000062068"
         },
       	{
           "ec": 1.4,
           "local_time": "2020-08-09 10:49:13.964827 +0900 KST m=+0.000117172"
         },
       	...
       ]
     }
     ```
     
     - local_time: (timestamp) 해당 기록의 등록 시간
   
   
   
6. /RegisterDevice

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

   

7. /RegisterRecipe

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

