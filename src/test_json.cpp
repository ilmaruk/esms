#include <iostream>
#include <string>

#include "nlohmann/json.hpp"

using json = nlohmann::json;

namespace ns
{
    struct person
    {
        std::string name;
        std::string address;
        int age;
    };

    void to_json(json &j, const person &p)
    {
        j = json{{"name", p.name}, {"address", p.address}, {"age", p.age}};
    }

    void from_json(const json &j, person &p)
    {
        j.at("name").get_to(p.name);
        j.at("address").get_to(p.address);
        j.at("age").get_to(p.age);
    }
}

json serialise(ns::person p)
{
    json j = p;
    std::cout << j << std::endl;
    return j;
}

int main()
{
    ns::person p{"Ned Flanders", "744 Evergreen Terrace, 60"};
    json j = serialise(p);

    auto p2 = j.template get<ns::person>();
    serialise(p2);

    return 0;
}