(ns adventofcode2016.core
  (:require [adventofcode2016.day05 :as day05])
  (:use [adventofcode2016.common])
  (:require [clojure.java.io :as io])
  (:use [clojure.pprint]))

(defn main
  []
  (-> ;"day05.txt"
      ;slurp-resource
      "done"
      day05/solve-part-2
      pprint))
