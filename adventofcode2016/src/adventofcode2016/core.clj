(ns adventofcode2016.core
  (:require [adventofcode2016.day02 :as day02])
  (:require [clojure.java.io :as io])
  (:use [clojure.pprint]))

(defn main
  []
  (-> "day02.txt"
      io/resource
      io/file
      slurp
      day02/solve-part-1
      pprint))
