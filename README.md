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
   8. /CheckDeviceOverlap



* Server domain

  Public DNS(IPv4): `https://asia-northeast1-superfarmers.cloudfunctions.net`



1. /Insert

   데이터베이스에 센서 데이터를 저장합니다.

   데이터 검수 과정 중 이상값이 감지되면 **abnormal** 컬렉션과 **desired_status** 컬렉션의 값을 업데이트합니다.

   | method |  path   |                request                |       response       |
   | :----: | :-----: | :-----------------------------------: | :------------------: |
   | `POST` | /Insert | (JSON) uuid를 포함한 센서 데이터 정보 | (string) 에러 메세지 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89779826-294a8080-db4b-11ea-8ea4-0059fad205f9.png)

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

     

     데이터에 이상값이 감지되면 **abnormal** 컬렉션엔 다음과 같은 정보가 기록됩니다.

     ![sample](https://user-images.githubusercontent.com/29545214/89992407-32109300-dcc0-11ea-84c0-bd85ede9a16d.png)

     - uuid: (string) 아두이노 기기의 고유 번호
     - errors: (Array&lt;number&gt;) 에러 코드
- time: (number) 에러 발생 시각 (UTC+0 기준)
     

     
다음은 `errors` 필드의 에러 코드들에 대한 명세입니다.
     
![sample](https://user-images.githubusercontent.com/29545214/89780645-d07be780-db4c-11ea-8592-3ce34112cbad.png)
     
또한 **desired_status** 컬렉션의 document id가 uuid와 일치하는 문서를 다음과 같이 업데이트합니다.
     
![sample](https://user-images.githubusercontent.com/29545214/89779997-847c7300-db4b-11ea-808d-235a7047af89.png)
     
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

     ![sample](https://user-images.githubusercontent.com/29545214/89780081-ae359a00-db4b-11ea-9804-9f5b4240a5c7.png)

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

     

3. /Control

   모바일 앱에서 아두이노 기기를 원격 제어하는 데 사용됩니다.

   **desired_status** 컬렉션의 상태를 변화시키고 아두이노 기기에서 데이터베이스의 변화를 감지하면 사용자가 설정한대로 움직이게 됩니다.

   | method |   path   |       request        |       response       |
   | :----: | :------: | :------------------: | :------------------: |
   | `POST` | /Control | (JSON) uuid와 상태값 | (string) 에러 메세지 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89780153-d7eec100-db4b-11ea-8784-b9c9a1bc625d.png)

     - uuid: (string) 대상 아두이노 기기의 고유 번호
     - led: (boolean) LED를 on/off
     - fan: (boolean) 팬을 on/off

     

4. /DesiredStatus

   아두이노 기기에서 사용자의 설정을 리스닝할 때 사용됩니다.

   **desired_status** 컬렉션의 상태를 받아 희망하는 상태값을 파악하고, 기기가 이 결과대로 움직이게 해야합니다.

   | method |      path      |    request    |    response    |
   | :----: | :------------: | :-----------: | :------------: |
   | `GET`  | /DesiredStatus | (string) uuid | (JSON) 상태 값 |

   - Query string 예시

     `uuid=123e6b776f000c04`

   - Response body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89780207-f05edb80-db4b-11ea-8a7d-b790e2345427.png)

     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
     - pH_pump: (number) pH 펌프 조절 정도 (1: pH 증가, 0: 유지, -1: pH 감소)
     - ec_pump: (number) ec 펌프 조절 정도 (1: ec 증가, 0: 유지)

     

5. /Records

   특정 필드의 최근 기록들을 불러옵니다.

   | method |   path   |              request              |      response      |
   | :----: | :------: | :-------------------------------: | :----------------: |
   | `POST` | /Records | (JSON) uuid와 필드명, 불러올 시간 | (JSON) 조회한 기록 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89780259-079dc900-db4c-11ea-9c8f-8ee4505c0347.png)

     - uuid: (string) 조회할 기기의 고유번호
     - key: (string) 조회할 필드명
     - time: (number) 조회할 기간(초 단위)

   

   - Response body 예시

     ​	![sample](https://user-images.githubusercontent.com/29545214/89780287-1edcb680-db4c-11ea-871f-73f91a7e1ed4.png)
     
     - local_time: (timestamp) 해당 기록의 등록 시간

   

6. /RegisterDevice

   새로운 기기 정보를 등록합니다.

   | method |      path       |               request               |       response       |
   | :----: | :-------------: | :---------------------------------: | :------------------: |
   | `POST` | /RegisterDevice | (JSON) 사용자 정보와 기기 고유 번호 | (string) 에러 메세지 |

   - Request body 예시

     ​	![sample](https://user-images.githubusercontent.com/29545214/89992605-756b0180-dcc0-11ea-9d6e-b1f46188bbe8.png)

     - uid: (string) 사용자 uid
     - uuid: (string) 기기 고유 번호

   

7. /RegisterRecipe

   작물 재배 레시피를 등록합니다.

   만약 같은 기기에 기존 레시피가 있었다면 덮어씁니다.

   | method |      path       | request |       response       |
   | :----: | :-------------: | :-----: | :------------------: |
   | `POST` | /RegisterRecipe | (JSON)  | (string) 에러 메세지 |

   - Request body 예시

     ​	![sample](https://user-images.githubusercontent.com/29545214/89992893-e7434b00-dcc0-11ea-8819-660b54f355f7.png)
     
     - uid: (string) 사용자 uid
     - uuid: (string) 기기 고유번호
     - crop: (string) 작물명
     - farm_name: (string) 농장명
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



8. /CheckDeviceOverlap

   기기 중복 여부를 검사합니다.

   등록돼있지 않은 기기라면 `404 Not Found` 에러를 반환합니다.

   사용중인 기기라면 `403 Forbidden` 에러를 반환합니다.

   등록 가능한 기기라면 `200 OK`를 반환합니다.
   
   | method |        path         |   request   |       response       |
| :----: | :-----------------: | :---------: | :------------------: |
   | `POST` | /CheckDeviceOverlap | (JSON) uuid | (string) 에러 메세지 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89897823-9842da80-dc1a-11ea-8c33-1d87513e8362.png)
   
     - uuid: (string) 기기 고유 번호

