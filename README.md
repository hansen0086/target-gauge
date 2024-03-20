# target-gauge

A prometheus exporter allows user to easily create metric and set its value by http request.

# build

    $docker build -t target-gauge .  

## run

    $docker run -p 127.0.0.1:80:8080/tcp target-gauge

## create new metric

    $curl -X GET "http://127.0.0.1:80/create?metric_name=my_gauge"

    # HELP craete a new metric named my_gauge  
    # TYPE target_gauge gauge
    my_gauge{target="any"} 0
## get metric

    $curl -X GET "http://127.0.0.1:80/metrics"

    # HELP target_gauge set value to /update to change value
    # TYPE target_gauge gauge
    target_gauge{target="any"} 0

## change metric value and get new value

    $curl -X GET "http://127.0.0.1:80/update?value=120.001&metric_name=my_gauge"
    updated%
    $curl -X GET http://127.0.0.1:80/metrics
    my_gauge{target="any"} 120.001

## delete metric

    $curl -X GET "http://127.0.0.1:80/delete?metric_name=my_gauge"





