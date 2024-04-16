- Data Transfer Object (DTO)
    ORM -> hibernate






class User{
    private int id;
    private String username;
    private String password;
    private String email;
    private String imgUrl;
    private Role role;
}

@ModelAttribute



DTO -> request dto

login -> email, password
signup -> username, email, password, imgUrl

LoginDto
SignUpDto

Dto -> response dto
profile -> username, email, imgUrl





- Null-pointer exception

    public User getUserById(int id){
        User user = new User();

        string q = "select * from";
        


        return user;
    }


    User user = UserDao.getUserById(id);

    if(user!=null && user.email == null){

    }



- equals() for String objects

1 == 2 ,  equals()

public void test(String username){
    if ("mg mg".equals(username)){

    }
}


- Don't use @Autowired

- REST API

- Git
