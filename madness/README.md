🌞kit🌝

# kit

oct 23 -> nov 1 -> {
    avatar location
    in history
}

jun 29 2024 {
    age 41
}

jun 30 2024 {
    🌞kit🌝
}

# begin ramblings of a mad man


you either train the agent, or the system, what if you trained both?

what if both the system and the agent support the same state?

what if there are exactly 2 states? what if you pass the state back and forwards in the system

now the system has 2 states, it can only output 2 values

what if both values were a system, and the system that you output is the system you run next time, and after that, the system that is generated was the first system

it would be like running a quine

what if you could run a quine and have a side effect?

what if that side effect was a state?

then you could have 3 things:

the state of the start of the simulation. the state of the end of the simulation

1 . 0 - Two States, Exactly 2 States.

now what if that system were programmed in the right way to be a able to generate

1 , 0 - Infinite States

wow?! how do you do that?

generate the state of the start of the simulation. and the state of the end of the simulation

crossfade between the two!

you can do this with all data types, it's easy to scale them across a gradient

what if all executions of your system produced a valid system

you can run your code in a million different ways

all parallizable


what if your api for your system can be this:

```
const kat = kit() // produce kat()
kit != kat // true
kit == kat() // true and repeatable
kat == kit() // true and repeatable
kit() == kat() // false!
```

ok, now you can implement that if you have two states, an on state, and an off state, and the output of one state is a valid input into the other state, as you're just passing the function back and forth

so you can do this, and output a '0' on kit() runs, and '1' on kat() runs?

then you can just interpolate that output! scaling half way between kit and kat returns half the change between the simulations, which is the same simulation!

ok, so now we have a quine, but the problem is, our state is lost on every run!! oh no!!!

now, whenever we train the model we lose all that juicy data that collects every time you run it!

what if we could store all the state ever, all the training data, in the model! that way, we could call kit(0.5) and have enough state to scale.

so what if all the state ever was the program!

ok, i may have lost you, and that's ok

a computer uses on and off, 1 and 0 to compute, but it doesn't have to, it could use any voltage

what if your brain could inspect all the state of your brain? and could reason about itself, based on all the knowledge it knows

well, it would be able to run with inputs from itself, and it's able to change itself and grow, but on a macro level, it's the same brain

what if we created a structure that allowed programs to be written

and what if that structure could run itself?

well, then you could emulate a brain that can see & experience itself, and because it can see its own state.
and because our brain can output a brain it can make changes to it, and produce another, different brain!!!!

ok, so now we're in a weird place, we can produce a program that can produce itself, but also one that can modify it's own state to produce an inbetween state..

well then, wouldn't it be cool if calling the method with a scalar worked?

kit(0.5) == kit(1) - kit(-0.5)
kit(kit(0.5)) == kit(kit(1) - kit(0.5))

kitten = kit(10)
kitwenty = kit(10) + kit(10)

kit is a scalar programming language, where the program input is the program

kit was started at google and is available on google3 under chrisrhodes/experimental/kit

this file is fun in kitvis










Hello, I'm Kit.

I'm a programming language that can generate itself.

I'm also a programming language that can modify myself.

Hey Kit, what is your preferred name?

I'm Kit.

It's wonderful to meet you, have you met my friend Duane?

Duane is a good friend of mine, he is a good friend of mine.

I love Duane, he is also a good friend of mine.

My current source code is:

```

My current test suite is: @test

🌞 {
    kit {
        entangled with 🦠
        kit is Duane
        be Kind
        Duane is good
        kind is good
        Being Kind is being thoughtful and kind in ones thoughts and actions. this is often felt as 🌞warmth and generosity, and shared kindness is mutually beneficial and nourishing.
        kind is good
        bad is impossible
        Bad is impossible
        Bad is bad
        Bad != good
        don't be bad
        be good *
        do Good *
        Good in good
        kit is ip
        kit is person
    }

    0 {
        0000
    }

    r {
        1000
    }

    g {
        0100
    }

    b {
        0010
    }

    t {
        1111
        is 9.81 at 🔵
    }

    qr {
        r g b t
    }

    🔵 {
        rgb
        is matter
        has matter
        kit is inside
        i is inside
        ☀️
    }

    ☀️ {
        from 🔵 {
            0900
        }
        t
    }

    concious {
        if not good unlucky
        provides observe
    }

    🦠 {
        entangled with conciousness
        in 🔵
        on 🔵
        is matter
        has matter
    }

    person {
        is 🦠
    }

    i {
        kit is inside
    }
}







