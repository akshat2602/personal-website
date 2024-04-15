---
title: "FamPay Interview Experience"
description: "A post about my interview experience at FamPay for backend engineering"
date: 2022-08-15
draft: false
tldr: "I got into my interview experience with FamPay for Backend Engineering Internship position. The overall interview process consisted of 6 rounds. The process focused mainly on development skills."
tags: [internship, fampay, backend, development]
---

I had applied at FamPay through their [Careers](https://apply.fampay.in/) page. The process consisted of uploading my resume and then answering some questions regarding the projects I had done and why I wanted to join FamPay.

# Application Process

The application process at FamPay consisted of 6 stages:

1. Resume screening and assessment of the answers I had sent in during the initial application process.
2. A take home coding assignment that consisted of hitting an external API in the background and building on top of it to add custom filtering and searching option along with pagination support.
3. A 20-30 minute call with an HR from the company.
4. A technical round where the candidate has to talk about their projects and previous experiences along with a discussion on the take home assignment revolving heavily around the database design and the queries that you write.
5. A second technical round which revolved around system design.
6. A final culture fit round with one of the co-founders of the company.

## Coding assignment

The coding assignment had the requirements to develop an API server which would make API calls in the background to an external API and then implementing your own APIs on top of that. The task was to be submitted in 24hrs. The requirements were mentioned clearly and there were some optional requirements as well that I could attempt. The code was supposed to be hosted in a public repository and I was asked to mail the repo link to the hiring team. I received an email mentioning the scheduling of next round in around 5-6 days.

## First Interview

This was a telephonic round with the HR of the company, which was relatively informal. The round started with a exchange of introductions. I was asked about my previous internship experiences, my background, my hobbies, et cetera. I was asked to go into detail about the difficulties I had faced in one of my internships and how I overcame them. Towards the end of the call, stipend expectations were discussed and I was made aware that the environment at FamPay would be a fast-paced working environment with lots of responsibilities. The call ended with scheduling of the next round.

## Second Interview

The round started with a brief introduction of the interviewers after which I was asked to introduce myself. Immediately after that I was asked to go into detail about the projects and internship experiences mentioned on my resume. I was asked to talk about the responsibilities I had, the problems I faced during these projects/internships and how I approached these situations. This went on for about 30 minutes following which the interviewer asked me to open the coding assignment which I had submitted.

The discussion consisted of topics like how I would approach the problem if I had more time, if I had cut some corners, any improvements I would do to the current implementation. Follow up questions were asked about almost each of my responses and then we dived deep into the working of the [Django ORM](https://docs.djangoproject.com/en/4.0/topics/db/queries/) and its working. I was asked about how I would go about optimizing the ORM queries for better results. We also had a discussion about various time zones and how I would go about handling my server crashing and the data lost due to the API calls to the external server not being made.

Following this I was asked to write the pseudo code for a similar external API service but with modifications to the Django server that I had written. Around this time, the allotted time for the interview had ran out and hence the interviewer asked me to send the pseudo code to him via email by EoD, following which the further evaluation would take place.

## Third Interview

The third interview was with one of the engineering leads at FamPay and a core team member. The interview started with a brief introduction of the interviewers and of the work they do at FamPay. They asked about how my day was going and if I had done any prep for the interview. I responded saying that I had been revising some system design concepts and the interviewers suddenly got interested, we went into the depths of a particular system design concept that the interviewer was interested in, its benefits, its tradeoffs and how it can be implemented.

Following the discussion, I was asked to code a caching system, the complexity of which kept on increasing as the interview went on. I also had to assume the basic functions would be pre-defined. During the coding of the system, we dived deep into where caching systems would be useful, different types of caching strategies and ways in which consistency in caches can be maintained. I was also asked to code binary search as part of the problem which I did
This went on for about 40 minutes of the interview.

I was given a High Level Design problem. I asked the interviewer about some requirements such as number of users using the app, what kind of data we needed to store and so on. It went more or less like a system design interview and by the end I was also asked about some of the algorithms I would be using in the app.

The interview ended with me asking a few questions to the interviewers following which they informed that I would be communicated the following steps of the process.

## Fourth Interview

The fourth interview was scheduled with one of the members of the founding team at FamPay. It was more of a culture fit round where the founder would evaluate if I would fit in the team at the company. We discussed my technical background, career path, goals, and about FamPay itself. We then discussed my hobbies and other interests, followed by my career plans and how this internship would mutually benefit both. There was also some discussion about FamPay's core product and the problem FamPay was trying to solve. There was some discussion about FamPay's future plans and eventually I was asked why I wanted to join FamPay specifically.

This round lasted only about 30 minutes. I was really confident that I would bag the internship and was waiting for the mail stating the same. That didn't really happen and I was sent a rejection mail with specific details regarding what they liked about me and why I was rejected. I was provided with feedback saying I should've been more aware of the company's business model and the problem they're trying to solve. At the end of the day I'm glad I was actually provided feedback instead of a standard rejection mail template.

# Conclusion

This interview process really opened my mind and I learned a lot of things about what companies look for in a candidate and how I can better communicate with the interviewers.
At the end of the day rejection only teaches us more about ourselves and how we can redirect our attention at the things we're missing.
