#### Smart Farmer Backend

![build status](https://github.com/CONCAT-internship/smartfarmer-backend/blob/master/assets/images/badge.svg)

<br>

**API specification**

0. Index
   1. /Insert
   3. /RecentStatus
   4. /Control
   5. /DesiredStatus
   6. /LookupByPeriod
   6. /LookupByNumber
   7. /RegisterDevice
   8. /RegisterRecipe
   9. /CheckDeviceOverlap
   10. /ProfileFarmer

<br>

* Server domain

  Public DNS(IPv4): `https://asia-northeast1-superfarmers.cloudfunctions.net`

<br>

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

   - Query param 예시
   
    `/RecentStatus?uuid=123e6b776f000c04`
       
       

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

   - Query param 예시

     `/DesiredStatus?uuid=123e6b776f000c04`

   - Response body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89780207-f05edb80-db4b-11ea-8a7d-b790e2345427.png)

     - led: (boolean) LED on/off
     - fan: (boolean) 팬 on/off
     - pH_pump: (number) pH 펌프 조절 정도 (1: pH 증가, 0: 유지, -1: pH 감소)
     - ec_pump: (number) ec 펌프 조절 정도 (1: ec 증가, 0: 유지)

     

5. /LookupByPeriod

   특정 기간 동안의 최근 기록들을 조회합니다.

   | method |      path       |              request              |      response      |
   | :----: | :-------------: | :-------------------------------: | :----------------: |
   | `POST` | /LookupByPeriod | (JSON) uuid와 필드명, 조회할 기간 | (JSON) 조회한 기록 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/90769444-36c1f080-e32b-11ea-87e6-97410050e2f9.png)

     - uuid: (string) 조회할 기기 uuid
     - key: (string) 조회할 필드명. key를 보내지 않을 경우 응답의 `data`필드에 모든 필드가 저장됩니다.
     - period: (number) 조회할 기간(초 단위)

   

   - Response body 예시 (key=ec인 경우)

     ![sample](https://user-images.githubusercontent.com/29545214/90769136-b7342180-e32a-11ea-9cd4-7406a6fb9d2c.png)
     
     - local_time: (timestamp) 해당 기록의 등록 시간 (오름차순으로 정렬돼있음)
     
     
     
   - Response body 예시 (key가 없을 경우)

     ![sample](https://user-images.githubusercontent.com/29545214/90769236-dcc12b00-e32a-11ea-92ce-26c03f39463c.png)

     

6. /LookupByNumber

   특정 갯수만큼의 최근 기록들을 조회합니다.

   | method |      path       |              request              |      response      |
   | :----: | :-------------: | :-------------------------------: | :----------------: |
   | `POST` | /LookupByNumber | (JSON) uuid와 필드명, 조회할 갯수 | (JSON) 조회한 기록 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/90769535-5c4efa00-e32b-11ea-8ebb-a4cc66cc0e7d.png)

     - uuid: (string) 조회할 기기 uuid
     - key: (string) 조회할 필드명. key를 보내지 않을 경우 응답의 `data`필드에 모든 필드가 저장됩니다.
     - number: (number) 조회할 갯수

     

   - Response body 예시 (key=ec인 경우)

     ![sample](https://user-images.githubusercontent.com/29545214/90769136-b7342180-e32a-11ea-9cd4-7406a6fb9d2c.png)

     - local_time: (timestamp) 해당 기록의 등록 시간 (오름차순으로 정렬돼있음)

     

   - Response body 예시 (key가 없을 경우)

     ![sample](https://user-images.githubusercontent.com/29545214/90769236-dcc12b00-e32a-11ea-92ce-26c03f39463c.png)

   

7. /RegisterDevice

   새로운 기기 정보를 등록합니다.

   | method |      path       |               request               |       response       |
   | :----: | :-------------: | :---------------------------------: | :------------------: |
   | `POST` | /RegisterDevice | (JSON) 사용자 정보와 기기 고유 번호 | (string) 에러 메세지 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89992605-756b0180-dcc0-11ea-9d6e-b1f46188bbe8.png)

     - uid: (string) 사용자 uid
     - uuid: (string) 기기 고유 번호

   

8. /RegisterRecipe

   작물 재배 레시피를 등록합니다.

   만약 같은 기기에 기존 레시피가 있었다면 덮어씁니다.

   | method |      path       | request |       response       |
   | :----: | :-------------: | :-----: | :------------------: |
   | `POST` | /RegisterRecipe | (JSON)  | (string) 에러 메세지 |

   - Request body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/89992893-e7434b00-dcc0-11ea-8819-660b54f355f7.png)
     
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



9. /CheckDeviceOverlap

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



10. /ProfileFarmer

    사용자의 프로필(닉네임과 등록된 스마트팜 정보)을 조회합니다.

    | method |      path      |       request       |       response       |
    | :----: | :------------: | :-----------------: | :------------------: |
    | `GET`  | /ProfileFarmer | (string) 사용자 uid | (JSON) 사용자 프로필 |

    - Query param 예시

      - `/ProfileFarmer?uid=Xecm2PHp7QNfCmb0MQOFdJdy5af2`

    - Response body 예시

      ![sample](https://user-images.githubusercontent.com/29545214/90770557-ec417380-e32c-11ea-81bd-34517721894f.png)

      - nickname: (string) 사용자 닉네임
      - device_uuid: (string) 기기 uuid
      - farm_name: (string) 농장명
