#include <iostream>


#include <cstdlib>
#include <string>
#include <vector>
#include <ctime>
#include <numeric>
#include <cmath>
#include <sstream>

#include <thread>
#include <chrono>
#include <mutex>
#include <thread>
#include <bit>

using namespace std;

static std::vector<unsigned char> bA1;
static std::vector<unsigned char> bA2;
static std::vector<unsigned char> bA3;
int total;

std::vector<unsigned char> writeArray(int size)
{
    //unsigned char buffer[1024];
    std::cout << "Starting to write Array" << endl;

    std::vector<unsigned char> byteArray;

    for (int i = 0; i <= size; i++)
    {
        byteArray.push_back(rand() % 2);
    }

    std::cout << "Done writing array" << endl;
    return byteArray;
}

std::vector<unsigned char> transfer(vector<unsigned char> input)
{
    std::vector<unsigned char> byteArrayTransfer;

    byteArrayTransfer = input;
    for (unsigned i = 1; i < input.size(); i++)
    {
        byteArrayTransfer.push_back(input[i]);
    }

    total += input.size();

    return byteArrayTransfer;
}


int main()
{
    static int size = 1024; //Setting the size of the arrays
    std::vector<std::thread> threads;
    std::cout << "Size of the Array is set to: " << size << endl<< endl;
    int cases = 1;
    int rounds = 0;


    auto startTime = chrono::steady_clock::now();

    std::thread t1([] {bA1 = writeArray(size); });
    std::thread t2([] {bA2 = writeArray(size); });
    std::thread t3([] {bA3 = writeArray(size); });




    t1.join();
    t2.join();
    t3.join();

    auto endTime = chrono::steady_clock::now();

    std::cout << "Done writing arrays. \nDuration: " << chrono::duration_cast<chrono::milliseconds>(endTime - startTime).count() << "ms" << endl;

    auto startTime1 = chrono::steady_clock::now();

    for (int i = 0; chrono::duration_cast<chrono::minutes>(chrono::steady_clock::now() - startTime1).count() < 1; i++)
    {

        switch (cases)
        {
        case 1:
        {
            std::thread t([] {bA2 = transfer(bA1); });
            threads.push_back(std::move(t));
        }
        break;

        case 2:
        {
            std::thread t([] {bA3 = transfer(bA2); });
            threads.push_back(std::move(t));
        }
        break;
        case 3:
        {
            std::thread t([] {bA1 = transfer(bA3); });
            threads.push_back(std::move(t));
        }
        break;
        }

        threads[i].join();

        if (cases == 1) cases = 2;
        if (cases == 2) cases = 3;
        if (cases == 3) cases = 1;

        rounds += 1;
    }




    auto endTime1 = chrono::steady_clock::now();

    std::cout << "\nDone transferring arrays." << endl << "Duration: " << chrono::duration_cast<chrono::seconds>(endTime1 - startTime1).count()
        << "s" << "\nSpeed: " << total / 1000000 * (chrono::duration_cast<chrono::seconds>(endTime1 - startTime1).count()) / 60
        << " mBit/s" << endl;
}