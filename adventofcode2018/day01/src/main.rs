use std::env;
use std::fs;
use std::collections::HashSet;

fn main() {
    // part_one();
    part_two();
}

fn part_two() {
    // let contents = "+1\n-2\n+3\n+1\n";
    let args: Vec<String> = env::args().collect();
    println!("{:?}", args);

    let contents = fs::read_to_string(&args[1])
       .expect("Something went wrong reading the file");

    let inputs: Vec<&str> = contents.split_whitespace().collect();

    let mut seen: HashSet<i32> = HashSet::new();

    let mut result: i32 = 0;
    let mut found = false;

    while !found {
        for s in &inputs {
            // println!("{}", s);

            if s.chars().count() > 1 {
                let op = &s[0..1];
                let x: i32 = s[1..].parse().unwrap();
                match op {
                    "+" => result += x,
                    "-" => result -= x,
                    _ => println!("failed on: {}", op),
                }

                if seen.insert(result) == false {
                    println!("done!");
                    found = true;
                    break;
                } else {
                    println!("{}", result);
                }
            }
        }
    }

    println!("Result {}", result);
}


fn part_one() {
    // let contents = "+1\n-2\n+3\n+1\n";
    let args: Vec<String> = env::args().collect();
    println!("{:?}", args);

    let contents = fs::read_to_string(&args[1])
       .expect("Something went wrong reading the file");

    let inputs: Vec<&str> = contents.split_whitespace().collect();

    let mut result: i32 = 0;

    for s in &inputs {
        println!("{}", s);

        if s.chars().count() > 1 {
            let op = &s[0..1];
            let x: i32 = s[1..].parse().unwrap();
            match op {
                "+" => result += x,
                "-" => result -= x,
                _ => println!("failed on: {}", op),
            }
            // println!("{}", result);
        }
    }

    println!("Result {}", result);
}